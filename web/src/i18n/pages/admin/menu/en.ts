export default {
  adminMenu: {
    add_menu:"Add menu",
    modify_menu:"Modify menu",
    Route :{
      id: "id",
      parent_id: "parent id",
      path: "path",
      name:  "name",
      component:  "component", // 对应前端文件路径
      title:  "title",
      isLink:  "isLink", //是否超链接菜单,开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
      isIframe:  "isIframe", //是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
      icon:  "icon",
    }
  },
};