package ch.qos.logback.classic.db.ext;

import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.UnsynchronizedAppenderBase;

public class MongoAppender extends UnsynchronizedAppenderBase<ILoggingEvent>{
	protected static final String STR_DATE_FORMAT = "yyyy-MM-dd HH:mm:ss";
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
	
	@Override
	public void append(ILoggingEvent eventObject) {
		mongo.append(eventObject);
	}

}
