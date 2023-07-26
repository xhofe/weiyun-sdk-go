package weiyunsdkgo

import "errors"

var (
	ErrCode403 = errors.New("http code 403")

	ErrTokenExpiration  = errors.New("token expiration")
	ErrTokenIsNil       = errors.New("token is nil")
	ErrCookieExpiration = errors.New("the login cookie is invalid, please login in again")
)

/*
 	190011: "无效的QQ号码",
    190012: "无效的命令字",
    190013: "请求参数错误",
    190014: "客户端主动取消,如关闭连接",
    190020: "组cmem包错误",
    190021: "解包cmem包错误",
    190030: "组ptlogin包失败",
    190031: "组pb协议包失败",
    190032: "解析pb协议包失败",
    190033: "解析http协议失败",
    190034: "解析json协议失败",
    190035: "解析xml协议失败",
    190036: "http状态码非200",
    190039: "无效的appid",
    190040: "UIN在黑名单中",
    190041: "Server内部错误",
    190042: "后端服务器超时",
    190043: "后端服务器进程不存在",
    190044: "解析后端回包失败",
    190045: "获取L5路由失败",
    190046: "服务器组包失败",
    190047: "严重错误，必须要引起重视",
    190048: "无效的APPID",
    190049: "可能违反互联网法律法规或腾讯服务协议",
    190050: "会话被强制下线",
    190051: "验证登录态失败",
    190052: "用户不在白名单中",
    190053: "用户在黑名单中",
    190054: "访问超过频率限制",
    190055: "服务器临时不可用",
    190056: "cmem key不存在",
    190057: "cmem key过期",
    190058: "cmem 没有数据",
    190059: "cmem 设置时cas不匹配",
    190060: "cmem 数据有误",
    190061: "无效的签名类型:请求身份验证凭证类型",
    190062: "解签名失败",
    190063: "解密数据失败",
    190064: "批量操作条目超上限",
    190065: "st签名过期，需要终端去换取新的Key",
    190066: "终端在同步的过程中，需要从头进行一次全量列表拉取",
    190067: "敏感文字",
    190071: "链接被对端关闭",
    190072: "策略限制",
    190201: "没有JSON头",
    190202: "没有JSON体",
    190203: "缺少必要参数",
    190204: "参数值类型不正确",
    199001: "回调callback参数异常",
    199002: "op_source参数有误",
    199003: "dir_key长度无效",
    199004: "文件sha长度无效",
    199005: "文件md5长度无效",
    199006: "",
    199007: "日志时间格式无效",
    199008: "域名不对",
    199009: "referer有问题",
    199010: "token有误",
    199011: "fileid长度无效",
    199012: "某参数超过配置限制",
    199013: "下载校验失败",
    199014: "用户请求信息非法",
    1e3: "服务器出错",
    1013: "存储平台系统繁忙",
    1015: "在存储平台创建用户失败",
    1016: "存储平台不存在该用户",
    1018: "要拉取的目录列表已经是最新的",
    1019: "目录不存在",
    1020: "文件不存在",
    1021: "目录ID已经使用",
    1022: "文件已传完",
    1026: "父目录不存",
    1027: "不允许在根目录下上传文件",
    1028: "目录或者文件数超过总限制",
    1029: "单个文件大小超限",
    1051: "重名错误",
    1052: "下载未完成上传的文件",
    1053: "当前上传的文件超过可用空间大小",
    1054: "不允许删除系统目录",
    1055: "不允许移动系统目录",
    1056: "该文件不可移动",
    1057: "续传时源文件已经发生改变",
    1058: "删除文件版本冲突",
    1059: "覆盖文件版本冲突",
    1060: "禁止查询根目录",
    1061: "禁止修改根目录属性",
    1062: "禁止删除根目录",
    1063: "不能删除非空目录",
    1064: "禁止拷贝未上传完成文件",
    1065: "不允许修改系统目录",
    1066: "原始外链url参数太长，超过了1022字节",
    1067: "短URL服务错误",
    1068: "短URL服务来源字段错误",
    1069: "短URL服务会数据包大小校验失败",
    1070: "生成外链文件大小不符合规则",
    1073: "外链失效，下载次数已超过限制",
    1074: "黑名单校验失败, 其它原因",
    1075: "黑名单校验失败，没有找到sha",
    1076: "非法文件，文件在黑名单中",
    1080: "名字太长",
    1081: "GET_APP_INFO时带的错误的source值",
    1082: "修改目录时间戳出错",
    1083: "目录或者文件数超单个目录限制",
    1084: "生成vaskey失败",
    1085: "批量操作不能为空",
    1086: "批量操作条目超上限",
    1088: "文件名目录名无效",
    1090: "无效的MD5",
    1091: "转存的文件未完成上传",
    1092: "转存的文件名无效编码",
    1093: "无效的业务ID",
    1094: "读取转存文件失败",
    1095: "转存文件已过期",
    1096: "设置flag失败",
    1097: "ftn preuploadblob解码失败",
    1098: "请求体中的业务号与业务blob中的业务号不一致",
    1099: "非法的目标业务号",
    1100: "微云preuploadblob解码失败",
    1101: "非法的文件前10M MD5",
    1102: "asn编码失败",
    1103: "存储存在此用户",
    1110: "转存到微云的文件名中含有非法字符",
    1111: "源、目的目录相同目录，不能移动文件",
    1112: "不允许文件或目录移动到根目录下",
    1113: "不允许文件复制到根目录下",
    1114: "移动索引不一致，存储需要修复",
    1115: "删除文件并发冲突,可以重试解决",
    1116: "不允许用户在根目录下创建目录",
    1117: "批量下载中某个目录或文件不存在",
    1118: "认证签名无效",
    1119: "目的父目录不存在",
    1120: "目的父父目录不存在",
    1121: "源父目录不存在",
    1122: "目录文件修改名称时，源目的相同",
    1123: "不允许在根目录下创建目录",
    1124: "访问旁路系统出错",
    1125: "黑名单",
    1126: "非秒传文件太大禁止上传",
    1301: "微云网盘用户不存在",
    1302: "QQ网盘用户不存在",
    1901: "独立密码签名已经超时，需要用户重新输入密码进行验证",
    1902: "独立密码验证失败",
    1903: "开通独立密码失败",
    1904: "删除独立密码失败",
    1905: "输入过于频繁",
    1906: "添加的独立密码和QQ密码相同",
    1908: "独立密码已经存在",
    1909: "修改密码失败",
    1910: "新老密码一样",
    1911: "不存在老密码，请用添加流程",
    1912: "策略限制",
    1913: "独立密码验证失败(密码错误)",
    1914: "失败次数过多,独立密码被锁定",
    1915: "认证签名无效",
    20418: "无效app_secret或access_token",
    20410: "微信access_token过期"
**/