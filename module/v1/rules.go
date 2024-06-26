package v1

type RuleData struct {
	Name string
	Type string
	Rule string
}

type Md5Data struct {
	Name   string
	Md5Str string
}

var RuleDatas = []RuleData{
	{"红帆OA", "body", "(iOffice.net)"},
	{"红帆-ioffice", "body", "(iOffice.net Hospital Edition|/iOffice/prg/welcome/welcomeShow.aspx?mode=login)"},
	{"红帆-HFOffice", "body", "(<title>HFOffice</title>|/config/pc-app-config.js)"},
	{"红帆-HFOffice", "icon_hash", "630899054"},
	{"泛微OA", "headers", "(ecology_JSessionid)"},
	{"泛微OA", "body", "(/spa/portal/public/index.js|wui/theme/ecology8/page/images/login/username_wev8.png|/wui/index.html#/?logintype=1)"},
	{"泛微OA", "icon_hash", "1578525679"},
	{"泛微云桥e-Bridge", "body", "(wx.weaver|e-Bridge)"},
	{"泛微E-Mobile", "body", "(Weaver E-mobile|移动管理平台-企业管理|weaver,e-mobile)"},
	{"泛微E-mobile", "headers", "(EMobileServer)"},
	{"致远OA", "body", "(/yyoa/|/seeyon/common/|/seeyon/USER-DATA/IMAGES/LOGIN/login.gif)"},
	{"致远M3 Server", "body", "(M3 Server|M3 Server V2.0)"},
	{"致远M1 Server", "body", "(M1-Server|M1-Server 已启动，您可以使用移动设备登录M1)"},
	{"致远M1 Server", "icon_hash", "3080850296"},
	{"金蝶OA", "body", "(EAS系统登录|金蝶国际软件集团有限公司)"},
	{"金蝶Apusic应用服务器", "body", "(欢迎使用Apusic应用服务器)"},
	{"金蝶云星空管理中心", "body", "(HTML5/content/themes/kdcss.min.css|/ClientBin/Kingdee.BOS.XPF.App.xap)"},
	{"金蝶政务GSiS", "body", "(/kdgs/script/kdgs.js|HTML5/content/themes/kdcss.min.css|/ClientBin/Kingdee.BOS.XPF.App.xap)"},
	{"万户OA", "body", "(/defaultroot/templates/template_system/common/css/|/defaultroot/scripts/|css/css_whir.css)"},
	{"万户ezOFFICE", "icon_hash", "2467445972"},
	{"蓝凌OA", "body", "(/scripts/jquery.landray.common.js|蓝凌OA|EKP|蓝凌软件|StylePath:\"/resource/style/default/\"|/resource/customization|sys/ui/extend/theme/default/style/profile.css|sys/ui/extend/theme/default/style/icon.css)"},
	{"蓝凌EIS智慧协同平台", "body", "(/scripts/jquery.landray.common.js)"},
	{"通达OA", "body", "(/static/images/tongda.ico|http://www.tongda2000.com|通达OA移动版|Office Anywhere|<link rel=\"shortcut icon\" href=\"/images/tongda.ico\" />)"},
	{"广联达Linkworks办公OA", "body", "(/Services/Identification/)"},
	{"广联达Linkworks办公OA", "icon_hash", "4005832539"},
	{"华天动力OA", "body", "(OAapp/WebObjects/OAapp.woa)"},
	{"协众OA", "body", "(scripts/cnoa.extra.js|admin@cnoa.cn|Powered by 协众OA|Powered by CNOA.CN)"},
	{"协众OA", "cookie", "(CNOAOASESSID)"},
	{"正方OA", "body", "(zfoausername)"},
	{"新点OA", "body", "(href=\"新点助手工具包.exe\")"},
	{"希尔OA", "body", "(/heeroa/login.do)"},
	{"中望OA", "body", "(/IMAGES/default/first/xtoa_logo.png|/app_qjuserinfo/qjuserinfoadd.jsp)"},
	{"海天OA", "body", "(HTVOS.js)"},
	{"信达OA", "body", "(http://www.xdoa.cn</a>)"},
	{"启莱OA", "body", "(window.location.href = MainPageFileName|Product.aspx?WebID=134)"},
	{"O2OA", "body", "(O2OA)"},
	{"O2OA", "icon_hash", "1367173865"},
	{"78OA", "body", "(78OA办公系统|http://www.78oa.com|license.78oa.com)"},
	{"易宝OA", "body", "(欢迎登录易宝OA系统)"},
	{"微宏OA", "body", "(wh/servlet/MainServer)"},
	{"佰鼎OA", "body", "(Skin2017)"},
	{"金合OA", "body", "(Jhsoft.Web.login)"},
	{"协达OA", "icon_hash", "-1850889691"},
	{"中腾OA", "body", "(zt_webframe)"},
	{"京伦OA", "body", "(京伦 OA系统)"},
	{"明致OA", "icon_hash", "1591287747"},
	{"信呼OA", "icon_hash", "1652488516"},
	{"同望iOA", "body", "(同望iOA协同办公平台)"},
	{"迪浪云OA", "icon_hash", "-179985729"},
	{"源天协同OA", "body", "(源天协同OA)"},
	{"然之协同OA", "body", "(然之协同OA)"},
	{"青年软件OA", "body", "(/Content/qrlib/ligerUI/skins/)"},
	{"天生创想知尚云OA", "body", "(天生创想 知尚云OA)"},
	{"九思OA协同办公系统", "body", "(/jsoa/login.jsp)"},
	{"泛普建筑工程施工OA", "body", "(/dwr/interface/LoginService.js)"},
	{"中软国际OA办公系统", "body", "(中软国际OA办公系统)"},
	{"RockOA协同办公OA系统", "body", "(RockOA协同办公OA系统)"},
	{"慧点科技OA协同办公系统", "body", "(/dojo/smartdot/css/dojo_smartdot.css)"},
	{"PHPOA开源协同OA办公系统", "body", "(PHPOA开源协同OA办公系统)"},
	{"用友-NC", "body", "(logo/images/ufida_nc.png|uclient.yonyou.com)"},
	{"用友-畅捷通OEM", "body", "(GNRemote.dll|Web_sc/login.gn)"},
	{"用友-畅捷通T+", "body", "(畅捷通 T+|/tplus/view/login.html)"},
	{"用友-UFIDA-NC", "body", "(nc.sfbase.applet.NCApplet.class)"},
	{"用友-ERP-NC", "body", "(/nc/servlet/nc.ui.iufo.login.Index)"},
	{"用友-Ufida", "body", "(/System/Login/Login.asp?AppID=)"},
	{"用友-U8", "body", "(getFirstU8Accid)"},
	{"用友-GRP-U8", "body", "(用友GRP-U8行政事业内控管理软件|用友GRP-U8(财务系统))"},
	{"用友-GRP-U8", "icon_hash", "3995446927"},
	{"用友-NC Cloud", "body", "(platform/pub/welcome.do|src=\"../../../../platform/resource/yonyou-yyy.js)"},
	{"用友-时空KSOA", "body", "(/images_index/productKSOA.jpg)"},
	{"用友-商战实践平台", "body", "(Login_Owner|Login_Main_BG)"},
	{"用友-移动系统管理系统", "body", "(<span></span>移动系统管理|../js/appActionLogCtr.js)"},
	{"用友-BIP数据应用服务", "icon_hash", "3263329630"},
	{"用友-IUFO", "body", "(iufo/web/css/menu.css)"},
	{"东华医疗协同办公系统", "body", "(extcomponent/security/login.jsp|login_dhc)"},
	{"商混ERP系统", "body", "(商混ERP系统|Microsoft Visual Studio .NET 7.1)"},
	{"企望制造ERP系统", "body", "(上海企望信息科技  技术支持|企望制造ERP系统)"},
	{"企望制造ERP系统", "icon_hash", "2054680991"},
	{"云时空社会化商业ERP系统", "title", "(云时空社会化商业ERP系统)"},
	{"云时空社会化商业ERP系统", "body", "(img src=\"/static/img/yunlogo_home.png)"},
	{"富通天下外贸ERP", "title", "(用户登录_富通天下外贸ERP)"},
	{"华夏ERP", "headers", "(华夏ERP)"},
	{"华夏ERP", "body", "(jshERP-boot)"},
	{"明源云ERP", "icon_hash", "-626335361"},
	{"明源云ERP", "title", "(明源云ERP)"},
	{"明源云ERP", "body", "(PubPlatform|明源云|LoginHandler)"},
	{"任我行-管家婆分销ERP", "body", "(分销ERP|txt_DbName|DownloadDriver.asp)"},
	{"若依管理系统", "body", "(ruoyi/login.js|ruoyi/js/ry-ui.js)"},
	{"若依管理系统", "icon_hash", "-1231872293"},
	{"宏景eHR人力资源信息管理系统", "body", "(人力资源信息管理系统|hrlogon)"},
	{"宏景HCM管理系统", "body", "(人力与人才信息管理系统)"},
	{"亿华人力资源管理系统", "body", "(<meta name=\"description\" content=\"亿华人力资源管理系统\" />)"},
	{"亿华人力资源管理系统", "icon_hash", "1429826889"},
	{"浙大恩特客户资源管理系统", "body", "(<title>欢迎使用浙大恩特客户资源管理系统</title>)"},
	{"时空智友企业信息管理系统存", "icon_hash", "3132337272"},
	{"易思无人值守智能物流系统", "body", "(Copyright 2005-2023 EOSINE Software Technology Company Limited all rights reserved|易思智能物流无人值守系统)"},
	{"易思无人值守智能物流系统", "icon_hash", "2007842404"},
	{"科荣AIO管理系统", "body", "(科荣控件|请下载<a href=\"/common/AIOPrint.exe\")"},
	{"科荣AIO管理系统", "icon_hash", "3656675207"},
	{"新开普前置服务管理平台", "title", "(掌上校园服务管理平台)"},
	{"新开普前置服务管理平台", "icon_hash", "3016838938"},
	{"契约锁电子签章系统", "body", "(电子签署平台管理后台|契约锁·电子签章|契约锁致力于为中大型组织提供电子签及印控一体化解决方案)"},
	{"契约锁电子签章系统", "icon_hash", "738823048"},
	{"契约锁电子签章系统", "icon_hash", "-2107882986"},
	{"紫光电子档案管理系统", "icon_hash", "3686048441"},
	{"紫光电子档案管理系统", "title", "(紫光档案管理系统——登录)"},
	{"禅道Zentao", "title", "(Welcome to use zentao|/theme/default/images/main/zt-logo.png)"},
	{"禅道Zentao", "headers", "(zentaosid)"},
	{"时空智友企业信息管理", "body", "(时空智友企业信息管理)"},
	{"章管家", "title", "(章管家登录-公章在外防私盖)"},
	{"畅捷通CRM", "title", "(畅捷CRM)"},
	{"畅捷通CRM", "icon_hash", "-1068428644"},
	{"智慧校园管理系统", "body", "(DC_Login/QYSignUp)"},
	{"蓝海卓越计费管理系统", "title", "(蓝海卓越计费管理系统)"},
	{"蓝海卓越计费管理系统", "body", "(蓝海卓越计费管理系统|星锐蓝海网络科技有限公司)"},
	{"大汉版通发布系统", "body", "(大汉版通发布系统|大汉网络)"},
	{"帆软报表 FineReport", "body", "(isSupportForgetPwd|ReportServer)"},
	{"帆软数据决策系统", "body", "(>数据决策系统|ReportServer?op)"},
	{"指挥调度管理平台", "icon_hash", "-1971504131"},
	{"指挥调度管理平台", "body", "(指挥调度管理平台)"},
	{"任我行CRM", "body", "(/Handlers/IdentifyingCode.ashx)"},
	{"东胜物流软件", "body", "(js/dhtmlxcombo_whp.js)"},
	{"智邦国际 企业管理软件", "body", "(update/exec.asp)"},
	{"银达汇智智慧综合管理平台", "body", "(福州银达云创信息科技有限公司|miniui/crypto/CodeManage.js)"},
	{"汉王人脸考勤管理系统", "body", "(汉王人脸考勤管理系统|/Content/image/hanvan.png|/Content/image/hvicon.ico)"},
	{"亿赛通-电子文档安全管理系统", "body", "(电子文档安全管理系统|/CDGServer3/index.jsp|/CDGServer3/SysConfig.jsp|/CDGServer3/help/getEditionInfo.jsp)"},
	{"ZKTeco考勤管理系统", "body", "(ZKTECO CO.,LTD.All rights reserved.)"},
	{"EnjoySCM供应链管理系统", "body", "(供应商网上服务厅)"},
	{"CoreMail", "body", "(coremail/common)"},
	{"亿邮电子邮件系统", "body", "(亿邮电子邮件系统|亿邮邮件整体解决方案)"},
	{"Spammark邮件信息安全网关", "body", "(/cgi-bin/spammark?empty=1)"},
	{"winwebmail", "body", "(WinWebMail Server|images/owin.css)"},
	{"atmail-WebMail", "cookie", "(atmail6)"},
	{"atmail-WebMail", "body", "(/index.php/mail/auth/processlogin|Powered by Atmail)"},
	{"TurboMail邮件系统", "body", "(TurboMail邮件系统|mailmain?intertype=ajax&type=getnodeid&uid=|Powered by TurboMail)"},
	{"飞企互联FE业务协作平台", "body", "(flyrise.stopBackspace.js)"},
	{"飞企互联FE业务协作平台", "icon_hash", "3903390150"},
	{"普元EOS", "icon_hash", "1204326113"},
	{"拓尔思TRS大数据平台", "body", "(Login page for TRS Media Asset Management System|<title>TRS媒资管理系统登录页面</title>)"},
	{"JeecgBoot企业级低代码平台", "body", "(var htmlRoot = document.getElementById|var theme = window.localStorage.getItem)"},
	{"海翔D8药业云平台", "title", "(登录海翔)"},
	{"YApi可视化接口管理平台", "body", "(YApi-高效、易用、功能强大的可视化接口管理平台|id=\"yapi\"|YApi|yapi接口管理)"},
	{"XXL-JOB", "body", "(分布式任务调度平台XXL-JOB)"},
	{"Smartbi", "body", "(smartbi.gcf.gcfutil)"},
	{"LiveBOS", "body", "(/react/browser/loginBackground.png|LiveBos)"},
	{"来客推商城系统", "title", "(来客推商城|来客推商城系统)"},
	{"来客推商城系统", "icon_hash", "1433892035"},
	{"狮子鱼CMS", "body", "(/seller.php?s=/Public/login)"},
	{"Joomla内容管理系统", "icon_hash", "1747282642"},
	{"PbootCMS", "headers", "(PbootSystem|PbootCMS)"},
	{"Dreamer CMS", "headers", "(dreamer-|dreamer-cms)"},
	{"Typecho", "headers", "(Typecho</a>|typecho|usr/themes)"},
	{"JEECMS", "body", "(/r/cms/www/red/js/common.js|/r/cms/www/red/js/indexshow.js|Powered by JEECMS|JEECMS|/jeeadmin/jeecms/index.do)"},
	{"Jspxcms", "body", "(- Powered by Jspxcms)"},
	{"ThinkPHP", "headers", "(ThinkPHP)"},
	{"ThinkPHP", "icon_hash", "1165838194"},
	{"ThinkPHP", "body", "(十年磨一剑-为API开发设计的高性能框架)"},
	{"ThinkPHP3", "body", "({ Fast & Simple OOP PHP Framework } -- [ WE CAN DO IT JUST THINK ])"},
	{"WordPress", "body", "(/wp-login.php?action=lostpassword|WordPress</title>)"},
	{"MetInfo CMS", "body", "(MetInfo|/skin/style/metinfo.css|/skin/style/metinfo-v2.css)"},
	{"EmpireCMS", "body", "(Powered by EmpireCMS)"},
	{"Druid", "body", "(druid.index|DruidDrivers|DruidVersion|Druid Stat Index)"},
	{"Discuz", "body", "(content=\"Discuz! X\")"},
	{"Drupal", "headers", "(drupal)"},
	{"Laravel", "headers", "(laravel_session)"},
	{"phpMyAdmin", "cookie", "(pma_lang|phpMyAdmin)"},
	{"phpMyAdmin", "body", "(/themes/pmahomme/img/logo_right.png)"},
	{"Emlog", "body", "(content=\"emlog\")"},
	{"Finecms", "body", "(content=\"FineCMS)"},
	{"74CMS", "body", "(74cms|qscms.root)"},
	{"大米CMS", "body", "(content=\"大米CMS|content=\"damicms)"},
	{"EduSoho教培系统", "title", "(Powered By EduSoho)"},
	{"EduSoho教培系统", "body", "(Powered by <a href=\"http://www.edusoho.com/\" target=\"_blank\">EduSoho|Powered By EduSoho)"},
	{"Ueditor", "body", "(ueditor.all.js|UE.getEditor)"},
	{"JeeSpringCloud", "body", "(com.jeespring.session.id)"},
	{"JeeSpringCloud", "icon_hash", "2446959922"},
	{"nocodb", "body", "(href=\"./_nuxt/nocodb)"},
	{"nocodb", "icon_hash", "2277371154"},
	{"Juniper Device Manager", "body", "(Juniper Web Device Manager)"},
	{"GeoServer", "body", "(GeoServer: 欢迎)"},
	{"GeoServer", "icon_hash", "375099679"},
	{"Apache Kylin", "body", "(url=kylin)"},
	{"Apache Dubbo", "headers", "(realm=\"dubbo\")"},
	{"Apache Druid", "headers", "(<title>Apache Druid</title>|content=\"Apache Druid console\")"},
	{"Apache Superset", "title", "(Superset)"},
	{"Apache Superset", "body", "(src=\"/static/assets/images/superset-logo-horiz.png\")"},
	{"Apache RocketMQ", "title", "(RocketMq)"},
	{"Apache RocketMQ", "body", "(RocketMq-console-ng|RocketMQ-Dashboard)"},
	{"Apache Struts2", "body", "(org.apache.struts2|Struts Problem Report|struts.devMode|struts-tags|There is no Action mapped for namespace)"},
	{"Apache Airflow", "title", "(Airflow - Login)"},
	{"Apache Activemq", "body", "(activemq_logo|Manage ActiveMQ broker)"},
	{"Apache Hadoop", "body", "(static/hadoop-st.png|/cluster/app/application)"},
	{"Adobe ColdFusion", "title", "(Error Occurred While Processing Request)"},
	{"JeecgBoot", "body", "(jeecg-boot)"},
	{"MinIO", "body", "(MinIO Console|MinIO Browser)"},
	{"Alibaba Nacos", "title", "(Nacos)"},
	{"Alibaba Nacos", "body", "(<title>Nacos</title>)"},
	{"Alibaba Druid", "body", "(click(druid.login.login|<title>druid monitor</title>|druid.common.buildHead)"},
	{"Metabase", "body", "(Metabase|/app/assets/img/apple-touch-icon.png)"},
	{"Nexus Repository Manager", "title", "(Nexus Repository Manager)"},
	{"Kibana", "title", "(Kibana)"},
	{"Weblogic", "body", "(/console/framework/skins/wlsconsole/images/login_WebLogic_branding.png|Welcome to Weblogic Application Server|<i>Hypertext Transfer Protocol -- HTTP/1.1</i>|<TITLE>Error 404--Not Found</TITLE>|Welcome to Weblogic Application Server|<title>Oracle WebLogic Server 管理控制台</title>)"},
	{"Weblogic", "headers", "(WebLogic)"},
	{"Influxdb", "headers", "(X-Influxdb)"},
	{"Shiro", "headers", "(rememberMe=|=deleteMe)"},
	{"Swagger UI", "body", "(/swagger-ui.css|swagger-ui-bundle.js|swagger-ui-standalone-preset.js|Swagger UI)"},
	{"Jboss", "body", "(Welcome to JBoss|jboss.css)"},
	{"Jboss", "headers", "(JBoss)"},
	{"GitLab", "body", "(assets/gitlab_logo)"},
	{"GitLab", "icon_hash", "1265477436"},
	{"GitLab", "icon_hash", "1278323681"},
	{"GitLab", "icon_hash", "516963061"},
	{"Zabbix", "icon_hash", "1045955894"},
	{"Zabbix", "icon_hash", "892542951"},
	{"Zabbix", "icon_hash", "2126845402"},
	{"Zabbix", "title", "(Zabbix)"},
	{"Spark", "icon_hash", "856048515"},
	{"RabbitMQ", "body", "(<title>RabbitMQ Management</title>)"},
	{"Tomcat", "body", "(/manager/status|/manager/html)"},
	{"SpringBoot", "icon_hash", "116323821"},
	{"SpringBoot", "body", "(Whitelabel Error Page)"},
	{"Atlassian Jira", "icon_hash", "981867722"},
	{"Atlassian Jira", "icon_hash", "552727997"},
	{"Atlassian Jira", "body", "(https://www.atlassian.com/software/jira|Atlassian Jira)"},
	{"Atlassian Confluence", "body", "(com.atlassian.confluence)"},
	{"Atlassian Confluence", "headers", "(X-Confluence)"},
	{"Atlassian Confluence", "icon_hash", "-305179312"},
	{"Atlassian Confluence", "icon_hash", "-1312806261"},
	{"Atlassian Bamboo", "icon_hash", "2221879069"},
	{"Outlook Web Application", "icon_hash", "1768726119"},
	{"F5 BIG-IP", "icon_hash", "-335242539"},
	{"F5 BIG-IP", "headers", "(BIG-IP® -|BIG-IP® Configuration Utility)"},
	{"Kubernetes", "body", "(Kubernetes Dashboard</title>|Kubernetes Enterprise Manager|Mirantis Kubernetes Engine|Kubernetes Resource Report)"},
	{"Jenkins", "body", "(登录 [Jenkins]|欢迎来到 Jenkins)"},
	{"Jenkins", "icon_hash", "81586312"},
	{"Harbor", "body", "(<title>Harbor</title>)"},
	{"Harbor", "cookie", "(harbor-lang)"},
	{"Grafana", "body", "(grafana-app|Grafana|<title>Grafana</title>)"},
	{"MeterSphere", "body", "(<title>MeterSphere</title>)"},
	{"鸿运主动安全云平台", "body", "(./open/webApi.html)"},
	{"鸿运主动安全云平台", "icon_hash", "3560704253"},
	{"锐捷网络", "body", "(static/img/title.ico|support.ruijie.com.cn|Ruijie - NBR|eg.login.loginBtn)"},
	{"锐捷BCR商业无线云网关", "body", "(<img src=\"/luci-static/ruijie/imgs/ruijlogo.jpg\")"},
	{"锐捷Ruijie", "body", "(4008 111 000)"},
	{"锐捷无线smartweb管理系统", "title", "(无线smartWeb--登录页面)"},
	{"汉塔科技上网行为流量管理系统", "body", "(汉塔科技 - 上网行为管理系统)"},
	{"汉塔科技上网行为流量管理系统", "icon_hash", "428165606"},
	{"联软IT安全运维管理系统", "icon_hash", "2163986542"},
	{"帕拉迪Core4A-UTM堡垒机", "body", "(帕拉迪Core4A-UTM)"},
	{"帕拉迪Core4A-UTM堡垒机", "icon_hash", "3928629595"},
	{"SDCMS神盾内容管理系统", "body", "(sdcms)"},
	{"傲盾信息安全管理系统", "body", "(aodun/aodun.js|傲盾软件)"},
	{"中新金盾安全管理与运维审计系统", "body", "(/fort_新版UI主线版本/WebRoot/pages/commons/meta.jsp)"},
	{"中新金盾信息安全管理系统", "body", "(中新金盾信息安全管理系统|中新网络信息安全股份有限公司)"},
	{"LiveBOS Manager管理控制平台", "body", "(LiveBOS控制台)"},
	{"LiveBOS Manager管理控制平台", "icon_hash", "3877875895"},
	{"汉得SRM云平台", "body", "(汉得SRM云平台)"},
	{"辰信景云终端安全管理系统", "body", "(辰信景云终端安全管理系统-SaaS版|辰信景云终端安全管理系统7.0)"},
	{"1Panel服务器运维管理面板", "icon_hash", "1965287190"},
	{"Milesight VPN", "title", "(MilesightVPN)"},
	{"Milesight VPN", "icon_hash", "612420834"},
	{"百卓Smart多业务安全网关", "body", "(writeCustomBgImg.jsp|baseajax)"},
	{"悦泰节能 智能数据网关", "body", "(/FWlib3/resources/css/xtheme-gray.cssz)"},
	{"惠尔顿e地通Socks5 VPN登录系统", "title", "(e地通Socks5 VPN登录系统)"},
	{"电信中兴ZXHN F450A网关", "title", "(ZXHN F450A)"},
	{"电信天翼网关F460", "title", "(F460)"},
	{"磊科Netcore", "title", "(磊科Netcore)"},
	{"西迪特Wi-Fi Web管理", "title", "(Wi-Fi Web管理)"},
	{"飞鱼星家用智能路由", "title", "(飞鱼星家用智能路由)"},
	{"飞鱼星上网行为管理", "body", "(css/R1Login.css|share.ti_username)"},
	{"网御星云WEB防护系统", "title", "(<title>网页防篡改系统 </title>|img src=\"images/loading.gif)"},
	{"网御VPN", "body", "(/vpn/common/js/leadsec.js|/vpn/user/common/custom/auth_home.css)"},
	{"浪潮Clusterengine V4.0集群管理平台", "body", "(TSCEV4.0 login)"},
	{"浪潮政务系统", "body", "(LangChao.ECGAP.OutPortal|OnlineQuery/QueryList.aspx)"},
	{"Kubepi", "title", "(Kubepi)"},
	{"Kubepi", "body", "(link rel=\"icon\" href=\"/kubepi/fav.png|link rel=\"icon\" href=\"/kubepi/link.ico)"},
	{"KubeOperator", "title", "(KubeOperator)"},
	{"ICEFLOW VPN Router系统", "title", "(ICEFLOW VPN Router)"},
	{"AVEVA InTouch安全网关", "body", "(InTouch Access Anywhere)"},
	{"nginxWebUI", "title", "(nginxWebUI)"},
	{"蜂网企业流控云路由器", "body", "(企业级流控云路由器|ifw8)"},
	{"TOTOLink路由器", "title", "(TOTOLINK)"},
	{"XiaoMi路由器", "body", "(<title>小米路由器</title>|<title>Redmi路由器</title>)"},
	{"Tenda路由器", "body", "(<title>Tenda)"},
	{"DrayTek Vigor路由器", "title", "(Vigor 2960)"},
	{"NetMizer日志管理系统", "body", "(NetMizer 日志管理系统)"},
	{"天玥运维安全网关", "body", "(css/fw/full.css|js/p/login.js)"},
	{"护卫神·主机大师", "icon_hash", "1188645141"},
	{"360网站安全检测", "body", "(webscan.360.cn/status/pai/hash)"},
	{"360网站卫士", "body", "(webscan.360.cn/status/pai/hash|wzws-waf-cgi|zhuji.360.cn/guard/firewall/stopattack.html)"},
	{"360网站卫士", "headers", "(360wzws|CWAP-waf|zhuji.360.cn|X-Safe-Firewall)"},
	{"360网神防火墙", "body", "(网神防火墙系统|resources/image/logo_header.png)"},
	{"360天堤新一代智慧防火墙", "body", "(360天堤)"},
	{"360天擎终端安全管理系统", "title", "(360天擎终端安全管理系统|360天擎)"},
	{"360天擎终端安全管理系统", "body", "(360新天擎|appid\":\"skylar6|已过期或者未授权，购买请联系4008-136-360|/task/index/detail?id={item.id})"},
	{"网神SecGate 3600防火墙", "body", "(网神SecGate|css/lsec/login.css|3600防火墙)"},
	{"NSFOCUS-防火墙", "body", "(NSFOCUS NF)"},
	{"NSFOCUS-防火墙", "headers", "(NSFocus)"},
	{"NSFOCUS-堡垒机", "body", "(/system/logo?id=large)"},
	{"NSFOCUS-下一代防火墙", "body", "(/stylesheet/iceye/images/logo/login_logo_nf_zh_CN.png)"},
	{"蓝盾防火墙", "body", "(Bluedon|default/js/act/login.js)"},
	{"Cisco SSLVPN", "body", "(/+CSCOE+/logon.html)"},
	{"Cisco CX20", "body", "(CISCO-CX20)"},
	{"TP-LINK", "headers", "(TP-LINK)"},
	{"TP-Link 3600 DD-WRT", "body", "(TP-Link 3600 DD-WRT)"},
	{"H3C 路由器", "body", "(/wnm/ssl/web/frame/login.html)"},
	{"H3C Web网管", "body", "(Web网管用户登录|china_logo.jpg)"},
	{"H3C ER2100n", "body", "(ER2100n系统管理)"},
	{"H3C ER8300G2", "body", "(ER8300G2系统管理)"},
	{"H3C ER5200G2", "body", "(ER5200G2系统管理)"},
	{"H3C ER6300", "body", "(ER6300系统管理)"},
	{"H3C ER6300G2", "body", "(ER6300G2系统管理)"},
	{"H3C ER3260", "body", "(ER3260系统管理)"},
	{"H3C ER3108G", "body", "(ER3108G系统管理)"},
	{"H3C ER3100", "body", "(ER3100系统管理)"},
	{"H3C SecBlade FireWall", "body", "(js/MulPlatAPI.js)"},
	{"H3C ER3108GW", "body", "(ER3108GW系统管理)"},
	{"H3C ER3260G2", "body", "(ER3260G2系统管理)"},
	{"H3C ICG1000", "body", "(ICG1000系统管理)"},
	{"H3C ER5200", "body", "(ER5200系统管理)"},
	{"华为（HUAWEI）SRG3250", "body", "(HUAWEI SRG3250)"},
	{"华为（HUAWEI）Secoway", "body", "(Secoway)"},
	{"华为（HUAWEI）SRG1220", "body", "(HUAWEI SRG1220)"},
	{"华为（HUAWEI）安全设备", "body", "(sweb-lib/resource/)"},
	{"华为（HUAWEI）USG", "body", "(UI_component/commonDefine/UI_regex_define.js)"},
	{"华为（HUAWEI）ASG2100", "body", "(HUAWEI ASG2100)"},
	{"SANGFOR-WAF", "body", "(commonFunction.js)"},
	{"SANGFOR-应用交付报表系统", "body", "(/reportCenter/index.php?cls_mode=cluster_mode_others)"},
	{"SANGFOR-行为感知系统", "body", "(isHighPerformance : !!SFIsHighPerformance,)"},
	{"SANGFOR-防火墙类产品", "body", "(SANGFOR FW)"},
	{"SANGFOR-SSL-VPN", "body", "(/por/login_psw.csp|loginPageSP/loginPrivacy.js)"},
	{"SANGFOR-EDR", "icon_hash", "1307354852"},
	{"SANGFOR-EDR", "body", "(/ui/static/img/title.1ac5aed.png)"},
	{"SANGFOR-上网行为管理系统", "body", "(login/gg/left_01.jpg)"},
	{"SANGFOR-上网行为管理系统", "body", "(utccjfaewjb = function|WRFWWCSFBXMIGKRKHXFJ)"},
	{"SANGFOR-应用交付管理系统", "icon_hash", "(http://www.sangfor.com.cn/product/it-yun-ad.html)"},
	{"天融信 VPN", "headers", "(topsecsvportalname)"},
	{"天融信防火墙", "body", "(TOPSEC|image/aaa.png)"},
	{"天融信TopAPP负载均衡系统", "body", "(TopAPP负载均衡系统)"},
	{"天融信脆弱性扫描与管理系统", "body", "(/js/report/horizontalReportPanel.js)"},
	{"天融信网络审计系统", "body", "(onclick=dlg_download())"},
	{"天融信日志收集与分析系统", "body", "(天融信日志收集与分析系统)"},
	{"天融信上网行为管理系统", "body", "(dkey_activex_download.php|images/logo3.gif)"},
	{"启明星辰天清汉马USG防火墙", "body", "(天清汉马USG|/cgi-bin/webui?op=get_product_model)"},
	{"启明星辰4A统一安全管控平台", "body", "(cas/css/ace-part2.min.css)"},
	{"启明星辰防火墙", "body", "(/cgi-bin/webui?op=get_product_model)"},
	{"群晖 NAS", "body", "(Synology|DiskStation|webman/modules)"},
	{"宝塔", "body", "(app.bt.cn/static/app.png|安全入口校验失败|<title>入口校验失败</title>|href=\"http://www.bt.cn/bbs|站点创建成功|检查端口是否正确|面板系统后台|宝塔Linux面板)"},
	{"IBM-Lotus-Domino", "body", "(/mailjump.nsf|/domcfg.nsf|/names.nsf|/homepage.nsf)"},
	{"向日葵-远程控制", "body", "({\"success\":false,\"msg\":\"Verification failure\"})"},
	{"TELEPORT堡垒机", "body", "(/static/plugins/blur/background-blur.js)"},
	{"JumpServer堡垒机", "body", "(<title>JumpServer</title>)"},
	{"思福迪堡垒机", "title", "(Logbase运维安全管理系统)"},
	{"中远麒麟堡垒机", "body", "(controller=admin_index&action=login)"},
	{"安恒云堡垒机", "body", "(安恒云堡垒机)"},
	{"安恒明御安全网关", "title", "(明御安全网关)"},
	{"安恒明御安全网关", "body", "(<title>明御安全网关</title>)"},
	{"安恒数据大脑API网关", "body", "(mssp-fe)"},
	{"瑞友天翼-应用虚拟化系统", "body", "(dvLogin|DownLoad.XGI|realor.cn|href=\"static/css/font-awesome.min.css|瑞友信息技术资讯有限公司)"},
	{"VMware vSphere", "body", "(VMware vSphere|ID_VISDK)"},
	{"Vmware Vcenter", "title", "(+ ID_VC_Welcome +)"},
	{"Vmware Vcenter", "icon_hash", "-175177211"},
	{"VMware Workspace ONE Access", "body", "(VMware Workspace ONE Access)"},
	{"ManageEngine ADManager Plus", "body", "(Hashtable.js|ADManager)"},
	{"Portainer（Docker管理）", "body", "(portainer.updatePassword|portainer.init.admin)"},
	{"Gogs简易Git服务", "cookie", "(i_like_gogs)"},
	{"Gitea简易Git服务", "cookie", "(i_like_gitea)"},
	{"Yundun", "headers", "(YUNDUN)"},
	{"Yunsuo", "headers", "(yunsuo)"},
	{"Safedog", "body", "(404.safedog.cn/images/safedogsite/broswer_logo.jpg)"},
	{"Safedog", "headers", "(Safedog|WAF/2.0)"},
	{"AWS S3 Bucket", "body", "(InvalidBucketName)"},
	{"Citrix Access Gateway", "body", "(Citrix Access Gateway)"},
	{"金山终端安全管理系统", "body", "(iepngfix/iepngfix_tilebg.js)"},
	{"金睛云华高级威胁检测系统", "icon_hash", "1747722638"},
	{"网防G01", "icon_hash", "-968234332"},
	{"HIKVISION iVMS-8700综合安防管理平台", "body", "(/portal/conf/icon/favicon.ico|/portal/conf/icon/logo.png)"},
	{"HIKVISION 综合安防管理平台", "title", "(综合安防管理平台)"},
	{"HIKVISION 综合安防管理平台", "body", "(home/locationIndex.action|dist/jquery.js)"},
	{"HIKVISION 对讲广播系统", "icon_hash", "-1830859634"},
	{"HIKVISION 视频监控", "icon_hash", "999357577"},
	{"大华视频监控", "icon_hash", "2019488876"},
	{"大华视频监控", "icon_hash", "833190513"},
	{"大华视频监控", "body", "(class=\"login-logo-dahua\")"},
	{"大华DSS", "body", "(<meta http-equiv=\"refresh\" content=\"1;URL='/admin'\"/>)"},
	{"大华智慧园区综合管理平台", "body", "(/WPMS/asset/lib/json2.js|DSS助手|/WPMS/asset/img/black/|src=\"/WPMS/asset/common/js/jsencrypt.min.js\")"},
	{"大华城市安防监控系统平台管理", "body", "(attachment_downloadByUrlAtt.action)"},
	{"NVS3000综合视频监控平台", "body", "(NVS3000综合)"},
	{"OfficeWeb365", "body", "(请输入furl参数|9999-9999-51A8-719E-C8B2-558D-11E9))"},
	{"XETUX", "body", "(<title>@XETUX - XPOS / BackEnd</title>)"},
	{"NUUO 摄像头", "body", "(Network Video Recorder Login)"},
	{"安美数字酒店宽带运营系统", "body", "(<title>酒店宽带运营系统</title>|<a href='http://www.amttgroup.com/' target='_blank'><font color='#FFFFFF'>安美数字)"},
	{"安美数字酒店宽带运营系统", "icon_hash", "1259797304"},
	{"企业微信-私有版服务端", "body", "(/wework_admin/static/style/images/independent/mulit_logo/WeworkLogoBule_2x$b672f477.png)"},
	{"企业微信-私有版服务端", "icon_hash", "3588968498"},
	{"魅课OM视频会议系统", "title", "(OMeeting视频会议|OM视频会议|OM免费网络视频会议系统)"},
	{"飞视美视频会议系统", "title", "(飞视美视频会议系统)"},
	{"会捷通云视讯平台", "body", "(him/api/rest/v1.0/node/role|him.app)"},
	{"好视通视频会议平台", "body", "(深圳银澎云计算有限公司|itunes.apple.com/us/app/id549407870|hao-shi-tong-yun-hui-yi-yuan)"},
	{"霆智科技VA虚拟应用平台", "body", "(EAA益和应用接入系统)"},
	{"TVT NVMS-1000", "body", "(NVMS-1000)"},
	{"小鱼易连云视讯管理平台", "body", "(font_1957344_lqkodjqdbl.css|static_source/localcdn/webrtc/web/favicon.ico)"},
	{"C-Lodop打印服务系统", "body", "(/CLodopfuncs.js|www.c-lodop.com)"},
	{"打印机", "body", "(打印机|media/canon.gif)"},
	{"目录遍历", "body", "(Directory listing for /|Index of /|- /</title>)"},
}
