package ch.qos.logback.classic.db.ext;

import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.rolling.RollingFileAppender;

public class MongoAndFileAppender extends RollingFileAppender<ILoggingEvent>{
	/**数据源。细化到数据库*/
	protected MongoConnectionSource dataSource;
	/**数据源。细化到数据库*/
	protected BaseMongo mongo;

	/**获取数据源。细化到数据库*/
	public MongoConnectionSource getMongoConnectionSource() {
	    return dataSource;
	}
	
	/**初始化数据源。细化到数据库*/
	public void setMongoConnectionSource(MongoConnectionSource dataSource) {
	    this.dataSource = dataSource;
    	mongo = new BaseMongo(dataSource);
	}
	/**是否继续往文件里写入*/
	private boolean appendFile = true;
	
	@Override
	public void append(ILoggingEvent eventObject) {
		//调用父类，往日志文件写入
		if(isAppendFile())
			super.append(eventObject);
		mongo.append(eventObject);
	}

	public boolean isAppendFile() {
		return appendFile;
	}

	public void setAppendFile(boolean appendFile) {
		this.appendFile = appendFile;
	}

}
