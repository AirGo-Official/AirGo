export default {
  adminMenu: {
    add_menu:"新增菜单",
    modify_menu:"修改菜单",
    Route :{
      id: "id",
      parent_menu:"父级菜单",
      parent_id: "父级菜单id",
      remarks:"别名",
      path: "路由路径",
      name:  "路由名称",
      component:  "前端文件路径", // 对应前端文件路径
      title:  "标题",
      isLink:  "外链地址", //是否超链接菜单,开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
      isIframe:  "外链内嵌", //是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
      icon:  "图标",
    }
  },
};