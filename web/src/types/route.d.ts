//路由信息
declare interface Route {
    created_at: string;
    updated_at: string;
    id: number;
    parent_id: number;
    path: string;
    name: string;
    component: string; // 对应前端文件路径
    children: Route[];
    roles: [];
    meta: {
        title: string;
        isLink: string; //是否超链接菜单,开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
        isIframe: boolean; //是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
        isHide: boolean
        isKeepAlive: boolean;//是否缓存组件状态
        isAffix: boolean; //是否固定在 tagsView 栏上
        icon: string;
    }
}

// 路由列表 used
declare interface RoutesListState {
    routesList: Route[];
    isColumnsMenuHover: Boolean;
    isColumnsNavHover: Boolean;
}

//used routes tree
declare interface RoutesTree {
    route_id: number,
    title: string,
    children?: RoutesTree[],
}

// 路由缓存列表
declare interface KeepAliveNamesState {
    keepAliveNames: string[];
    cachedViews: string[];
}

// 后端返回原始路由(未处理时)
declare interface RequestOldRoutesState {
    requestOldRoutes: string[];
}

// TagsView 路由列表
declare interface TagsViewRoutesState<T = any> {
    tagsViewRoutes: T[];
    isTagsViewCurrenFull: Boolean;
}


