export default {
  adminMenu: {
    add_menu:"新增選單",
    modify_menu:"修改選單",
    Route :{
      id: "id",
      parent_menu:"父級選單",
      parent_id: "父級選單id",
      remarks:"别名",
      path: "路由路徑",
      name:  "路由名稱",
      component:  "前端檔案路徑", // 对应前端文件路径
      title:  "標題",
      isLink:  "外鏈地址", //是否超链接菜单,开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
      isIframe:  "外鏈內嵌", //是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
      icon:  "圖標",
    }
  },
};