// 布局配置
declare interface ThemeConfigState {
    themeConfig: {
        created_at: string;
        updated_at: string;
        id: number;
        [bg:string];
        isDrawer: boolean;
        primary: string;
        topBar: string;
        topBarColor: string;
        isTopBarColorGradual: boolean;
        menuBar: string;
        menuBarColor: string;
        menuBarActiveColor: string;
        isMenuBarColorGradual: boolean;
        columnsMenuBar: string;
        columnsMenuBarColor: string;
        isColumnsMenuBarColorGradual: boolean;
        isColumnsMenuHoverPreload: boolean;
        isCollapse: boolean;
        isUniqueOpened: boolean;
        isFixedHeader: boolean;
        isFixedHeaderChange: boolean;
        isClassicSplitMenu: boolean;
        isLockScreen: boolean;
        lockScreenTime: number;
        isShowLogo: boolean;
        logo_link: string; //logo链接
        isShowLogoChange: boolean;
        isBreadcrumb: boolean;
        isTagsview: boolean;
        isTagsviewIcon: boolean;
        isCacheTagsView: boolean;
        isSortableTagsView: boolean;
        isShareTagsView: boolean;
        isFooter: boolean;
        isGrayscale: boolean;
        isInvert: boolean;
        isIsDark: boolean;
        isWartermark: boolean;
        wartermarkText: string;
        tagsStyle: string;
        animation: string;
        columnsAsideStyle: string;
        columnsAsideLayout: string;
        layout: string;
        isRequestRoutes: boolean;
        globalTitle: string;
        globalViceTitle: string;
        globalViceTitleMsg: string;
        globalI18n: string;
        globalComponentSize: string;
    };
}

// aside
declare type AsideState = {
    menuList: RouteRecordRaw[];
    clientWidth: number;
};

// columnsAside
declare type ColumnsAsideState<T = any> = {
    columnsAsideList: T[];
    liIndex: number;
    liOldIndex: null | number;
    liHoverIndex: null | number;
    liOldPath: null | string;
    difference: number;
    routeSplit: string[];
};

// navBars breadcrumb
declare type BreadcrumbState<T = any> = {
    breadcrumbList: T[];
    routeSplit: string[];
    routeSplitFirst: string;
    routeSplitIndex: number;
};

// navBars search
declare type SearchState<T = any> = {
    isShowSearch: boolean;
    menuQuery: string;
    tagsViewList: T[];
};

// navBars tagsView
declare type TagsViewState<T = any> = {
    routeActive: string | T;
    routePath: string | unknown;
    dropdown: {
        x: string | number;
        y: string | number;
    };
    sortable: T;
    tagsRefsIndex: number;
    tagsViewList: T[];
    tagsViewRoutesList: T[];
};

// navBars parent
declare type ParentViewState<T = any> = {
    refreshRouterViewKey: string;
    iframeRefreshKey: string;
    keepAliveNameList: string[];
    iframeList: T[];
};

// navBars link
declare type LinkViewState = {
    title: string;
    isLink: string;
};
