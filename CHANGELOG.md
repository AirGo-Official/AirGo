# Changelog

## 2024.5.3
### #001
feat:在”首页“中添加了分栏、添加“概览”卡片，并可进行部分快捷操作
feat:在”首页“的订阅卡片中添加“⚡一键订阅”功能（目前已添加Clash Meta），订阅卡片按钮位置调整
feat:添加与更改了部分ico
feat:在”工单“页面添加了工单数量统计
feat:调整了“使用文档”页面的布局，并在阅读的页面中，为标题、副标题新增卡片、并添加了返回按钮
feat:更新vue相关依赖至最新版本（包括Element plus）
feat:增加i18n条目
refactor:修改了易支付跳转页面的样式
refactor:优化了顶栏、面包屑的样式移动端适配
refactor:“商店”页面优化、移动端适配
refactor:部分UI与组件优化
refactor:优化了工单页面的表格、排序进行调整
refactor:调整超时时间至10s
fix:修复了深色模式导致文章中的文字内容无法辨别的问题
fix:修复了深色模式中菜单背景透明的问题
pref:删除了部分无用参数

## 2024.5.
### #002
feat:增加静态资源 assets api 开关按钮，在前后端分离时，如果不想AirGo核心提供静态资源api，可关闭 assets api
feat:增加swagger api 开关按钮，具体使用请看 swagger api
feat:菜单样式更换
feat:面包屑样式优化
feat:ico样式进行了替换更新
feat:#41 订阅名称自定义，在部分客户端一键导入的时候有效果。例如：http://192.168.0.61:8899/api/public/sub/f347004b3b0645e18dc7da07ea5c6f92/AirGo%E7%89%9B%E9%80%BC，其中AirGo%E7%89%9B%E9%80%BC为订阅名称，可以任意修改。
feat:#42 Clash和Surge自定义分流规则
feat:前后分离时，前端可以对接多个api地址，以英文符号|分割。在index.html中修改，window.httpurl = "http://192.168.0.1:5555|http://192.168.0.2:6666|http://192.168.0.3:8888"
feat:登录页样式更改
refactor:购买、续费商品的样式调整
refactor:修复了部分i18n错误
refactor:进行多端的ui优化
fix:#39 修复分页查询问题
fix:修复了多处前端错误
pref:删除了部分无用参数
pref:删除了部分无效的i18n
pref:优化打开速度、替换了部分引用外部cdn的文件