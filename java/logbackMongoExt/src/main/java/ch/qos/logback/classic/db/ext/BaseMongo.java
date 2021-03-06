package ch.qos.logback.classic.db.ext;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.locks.ReentrantLock;

import org.bson.Document;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.mongodb.MongoClient;
import com.mongodb.MongoClientOptions;
import com.mongodb.MongoClientOptions.Builder;
import com.mongodb.MongoCredential;
import com.mongodb.ServerAddress;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;

import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.classic.spi.IThrowableProxy;
import ch.qos.logback.classic.spi.StackTraceElementProxy;

public class BaseMongo {
	protected static final Logger LOGGER = LoggerFactory.getLogger(BaseMongo.class);
	
	protected final ReentrantLock lock = new ReentrantLock(true);
	
	/** 插入错误次数 */
	protected int errorCount;
	
	/**数据源。细化到数据库*/
	protected MongoConnectionSource dataSource;
	
	/**mongo数据库*/
	protected MongoDatabase db;
	
	/**当前日志结合，用来做日志表分表的判断*/
	protected String currCollection;
	/**分表的时间格式化参数*/
	protected String collectionPattern;
	
	protected DateFormat df;
	
	public BaseMongo(MongoConnectionSource dataSource) {
		try {
			System.out.println(this.getClass().getName()+" mongo init start");
			initMongo(dataSource);
			System.out.println(this.getClass().getName()+" mongo init end");
		} catch(Throwable t) {
			t.printStackTrace();
			LOGGER.error("mongo初始化异常", t);
			//程序启动，如果mongo出错，强制退出
			System.exit(3);
		}
	}

	/**获取数据源。细化到数据库*/
	public MongoConnectionSource getMongoConnectionSource() {
	    return dataSource;
	}

	protected MongoClient mongo;
	protected MongoCollection<Document> table;
	
	
	public void append(ILoggingEvent eventObject) {
		try {
		Document doc = new Document();
		
		//可以用eventObject.getMDCPropertyMap()获取之前设置好的上下文信息。
	    //MDC.put(key, val)。可以放入session和用户id等。先设置mdc，避免mdc与日志系统属性覆盖
		Map<String, String> mdc = eventObject.getMDCPropertyMap();
	    if (mdc != null && !mdc.isEmpty()) {
	    	for(String key : mdc.keySet()) {
    			doc.append(key, mdc.get(key));
	    	}
	    }
	    
	    Date logDate = new Date(eventObject.getTimeStamp());
		doc.append("log_date", logDate);
	    doc.append("level", eventObject.getLevel().toString());
	    doc.append("logger", eventObject.getLoggerName());
	    doc.append("thread", eventObject.getThreadName());

	    //可以获取栈的信息
	    StackTraceElement first = eventObject.getCallerData()[0];
	    doc.append("className", first.getClassName());
	    doc.append("methodName", first.getMethodName());
	    doc.append("lineNum", first.getLineNumber());
	    
	    //格式化后的消息，eventObject.getMessage()是格式化之前的文本
	    doc.append("message", eventObject.getFormattedMessage());
	    
	    //是否存在异常
	    IThrowableProxy throwable = eventObject.getThrowableProxy();
		if(throwable!=null) {
			if(throwable.getMessage()!=null)
				doc.append("errorMessage", throwable.getMessage());
			doc.append("throwable", toMongoDocument(throwable));
	    }
		
		
		//日志表进行分表判断
    	if(collectionPattern!=null) {
    		setSubTable(logDate);
    	}

    	table.insertOne(doc);
        	
		} catch (Throwable t) {
			errorCount++;
			//TODO 需要通知运维及开发人员
			LOGGER.error("mongo保存异常,errorCount:"+errorCount, t);
		}
	}
	
	/**
	 * 
	 * @Title: setSubTable
	 * @Description: 设置分表
	 * @param @param logDate   
	 * @return void   
	 * @throws
	 */
	private void setSubTable(Date logDate) {
		lock.lock();
		try {
			String suffix = df.format(logDate);
			if(currCollection.endsWith(suffix)) {
				return;
			}
			
			currCollection = dataSource.getCollection()+suffix;
			table = db.getCollection(currCollection);
		} finally {
			lock.unlock();
		}
	}

	/**
	 * 获取异常信息
	 * @param throwable
	 * @return
	 */
	protected Document toMongoDocument(IThrowableProxy throwable) {
        Document throwableDoc = new Document();
        throwableDoc.append("class", throwable.getClassName());
        throwableDoc.append("message", throwable.getMessage());
        throwableDoc.append("stackTrace", toSteArray(throwable));
        if (throwable.getCause() != null) {
            throwableDoc.append("cause", toMongoDocument(throwable.getCause()));
        }
        return throwableDoc;
    }
	
	/**
	 * 获取异常栈
	 * @param throwableProxy
	 * @return
	 */
	protected String toSteArray(IThrowableProxy throwableProxy) {
		StringBuilder buf = new StringBuilder();
        final StackTraceElementProxy[] elementProxies = throwableProxy.getStackTraceElementProxyArray();
        final int totalFrames = elementProxies.length - throwableProxy.getCommonFrames();
        for (int i = 0; i < totalFrames; ++i)
        	buf.append(elementProxies[i].getStackTraceElement().toString()).append(";");
        return buf.toString();
    }
	
	/**
	 * 获取栈的信息
	 * @param callerData
	 * @return
	 */
	protected List<Document> toDocument(StackTraceElement[] callerData) {
        LinkedList<Document> list = new LinkedList<Document>();
        for (final StackTraceElement ste : callerData) {
        	Document document = new Document()
                    .append("file", ste.getFileName())
                    .append("class", ste.getClassName())
                    .append("method", ste.getMethodName())
                    .append("line", ste.getLineNumber())
                    .append("native", ste.isNativeMethod());
        	list.add(document);
        }
        return list;
    }
	
	/**初始化mongo*/
	private void initMongo(MongoConnectionSource dataSource) {
	    this.dataSource = dataSource;
	    
		//设置鉴权信息
		List<MongoCredential> credentialsList = new LinkedList<MongoCredential>();
		if(dataSource.getUsername()!=null&&dataSource.getPassword()!=null) {
			MongoCredential user = MongoCredential.createCredential(dataSource.getUsername(), dataSource.getDb(), dataSource.getPassword().toCharArray());
			credentialsList.add(user);
		}
		
		//设置参数
		Builder optBuilder = MongoClientOptions.builder();
		if(dataSource.getSocketTimeout()>0)
			optBuilder.socketTimeout(dataSource.getSocketTimeout());
		if(dataSource.getConnectTimeout()>0)
			optBuilder.connectTimeout(dataSource.getConnectTimeout());
		if(dataSource.getMaxWaitTime()>0)
			optBuilder.maxWaitTime(dataSource.getMaxWaitTime());
		optBuilder.socketKeepAlive(dataSource.isSocketKeepAlive());
		MongoClientOptions opt = optBuilder.build();
		
		//创建数据库客户端
		if(credentialsList.size()==0) {
			mongo = new MongoClient(new ServerAddress(dataSource.getUrl()), opt);
		} else {
			mongo = new MongoClient(new ServerAddress(dataSource.getUrl()), credentialsList, opt);
		}
		db = mongo.getDatabase(dataSource.getDb());
		
		//设置日志表
		currCollection = dataSource.getCollection();
		collectionPattern = dataSource.getCollectionPattern();
		if(collectionPattern!=null) {
			collectionPattern = collectionPattern.trim();
			if(collectionPattern.length()==0) {
				throw new IllegalArgumentException("collectionPattern is blank:"+collectionPattern);
			}
		}
		if(collectionPattern!=null) {
			df = new SimpleDateFormat(this.collectionPattern);
			setSubTable(new Date());
		} else {
			table = db.getCollection(currCollection);
		}
	}

}
