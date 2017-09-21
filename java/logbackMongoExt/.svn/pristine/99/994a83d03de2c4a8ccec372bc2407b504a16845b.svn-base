package ch.qos.logback.classic.db.ext;

/**
 * mongo数据库配置信息
 * @author jinlong
 *
 */
public class MongoConnectionSource {
	protected String username;
	protected String password;
	protected String url;
	protected String db;
	protected String collection;
	private String collectionPattern;
	protected int socketTimeout;
	protected int maxWaitTime;
	protected int connectTimeout;
	protected boolean socketKeepAlive;
	
	public String getUsername() {
		return username;
	}
	public void setUsername(String username) {
		this.username = username;
	}
	public String getPassword() {
		return password;
	}
	public void setPassword(String password) {
		this.password = password;
	}
	public String getUrl() {
		return url;
	}
	public void setUrl(String url) {
		this.url = url;
	}
	public String getDb() {
		return db;
	}
	public void setDb(String db) {
		this.db = db;
	}
	public String getCollection() {
		return collection;
	}
	public void setCollection(String collection) {
		this.collection = collection;
	}
	public int getSocketTimeout() {
		return socketTimeout;
	}
	public void setSocketTimeout(int socketTimeout) {
		this.socketTimeout = socketTimeout;
	}
	public boolean isSocketKeepAlive() {
		return socketKeepAlive;
	}
	public void setSocketKeepAlive(boolean socketKeepAlive) {
		this.socketKeepAlive = socketKeepAlive;
	}
	public int getMaxWaitTime() {
		return maxWaitTime;
	}
	public void setMaxWaitTime(int maxWaitTime) {
		this.maxWaitTime = maxWaitTime;
	}
	public int getConnectTimeout() {
		return connectTimeout;
	}
	public void setConnectTimeout(int connectTimeout) {
		this.connectTimeout = connectTimeout;
	}
	public String getCollectionPattern() {
		return collectionPattern;
	}
	public void setCollectionPattern(String collectionPattern) {
		this.collectionPattern = collectionPattern;
	}
	@Override
	public String toString() {
		return "MongoConnectionSource [username=" + username + ", url=" + url + ", db=" + db
				+ ", collection=" + collection + ", collectionPattern=" + collectionPattern + ", socketTimeout="
				+ socketTimeout + ", maxWaitTime=" + maxWaitTime + ", connectTimeout=" + connectTimeout
				+ ", socketKeepAlive=" + socketKeepAlive + "]";
	}
	
}
