#javac -d . gui/PubSubSQLGUI.java gui/MainForm.java
#jar cvmf gui/main_manifest pubsubsqlgui.jar PubSubSQLGUI.class PubSubSQLGUI1$.class MainForm.class gui/images/New.png
#java -jar pubsubsqlgui.jar
 
javac -d . gui/PubSubSQLGUI.java gui/MainForm.java
jar cvf pubsubsqlgui.jar PubSubSQLGUI.class MainForm.class gui/images/New.png
jar cfe pubsubsqlgui.jar PubSubSQLGUI PubSubSQLGUI.class
java -jar pubsubsqlgui.jar