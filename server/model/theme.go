package model

import "time"

type Theme struct {
	// 是否开启布局配置抽屉
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"   gorm:"primaryKey"`

	IsDrawer bool `json:"isDrawer" gorm:"default:false;comment:是否开启布局配置抽屉"`
	//全局主题
	Primary  string `json:"primary"  gorm:"default:#002EFC;comment:默认primary主题颜色"`
	IsIsDark bool   `json:"isIsDark" gorm:"default:false;comment:是否开启深色模式"`

	//顶栏设置
	TopBar               string `json:"topBar"               gorm:"default:#ffffff;comment:默认顶栏导航背景颜色"`
	TopBarColor          string `json:"topBarColor"          gorm:"default:#606266;comment:默认顶栏导航字体颜色"`
	IsTopBarColorGradual bool   `json:"isTopBarColorGradual" gorm:"default:false;comment:是否开启顶栏背景颜色渐变"`

	// 菜单设置
	MenuBar               string `json:"menuBar"               gorm:"default:#FFFFFF;comment:默认菜单导航背景颜色"`
	MenuBarColor          string `json:"menuBarColor"          gorm:"default:#000000;comment:默认菜单导航字体颜色"`
	MenuBarActiveColor    string `json:"menuBarActiveColor"    gorm:"default:'rgba(0, 0, 0, 0.2)';comment:默认菜单高亮背景色"`
	IsMenuBarColorGradual bool   `json:"isMenuBarColorGradual" gorm:"default:false;comment:是否开启菜单背景颜色渐变"`

	// 分栏设置
	ColumnsMenuBar               string `json:"columnsMenuBar"                  gorm:"default:#FFFFFF;comment:默认分栏菜单背景颜色"`
	ColumnsMenuBarColor          string `json:"columnsMenuBarColor"            gorm:"default:#000000;comment:默认分栏菜单字体颜色"`
	IsColumnsMenuBarColorGradual bool   `json:"isColumnsMenuBarColorGradual" gorm:"default:false;comment:是否开启分栏菜单背景颜色渐变"`
	IsColumnsMenuHoverPreload    bool   `json:"isColumnsMenuHoverPreload"     gorm:"default:false;comment:是否开启分栏菜单鼠标悬停预加载(预览菜单)"`

	//界面设置
	IsCollapse          bool  `json:"isCollapse"            gorm:"default:false;comment:是否开启菜单水平折叠效果"`
	IsUniqueOpened      bool  `json:"isUniqueOpened"       gorm:"default:true;comment:是否开启菜单手风琴效果"`
	IsFixedHeader       bool  `json:"isFixedHeader"        gorm:"default:true;comment:是否开启固定 Header"`
	IsFixedHeaderChange bool  `json:"isFixedHeaderChange" gorm:"default:false;comment:初始化变量，用于更新菜单 el-scrollbar 的高度，请勿删除"`
	IsClassicSplitMenu  bool  `json:"isClassicSplitMenu"  gorm:"default:false;comment:是否开启经典布局分割菜单（仅经典布局生效）"`
	IsLockScreen        bool  `json:"isLockScreen"         gorm:"default:false;comment:是否开启自动锁屏"`
	LockScreenTime      int64 `json:"lockScreenTime"       gorm:"default:30;comment:开启自动锁屏倒计时(s/秒)"`

	//界面显示
	IsShowLogo         bool   `json:"isShowLogo"         gorm:"default:true;comment:是否开启侧边栏 Logo"`
	LogoLink           string `json:"logo_link"          gorm:"default:https://telegraph-image.pages.dev/file/c48a2f45ebf102dd66131.png;comment:logo链接"`
	IsShowLogoChange   bool   `json:"isShowLogoChange"   gorm:"default:false;comment:初始化变量，用于 el-scrollbar 的高度更新，请勿删除"`
	IsBreadcrumb       bool   `json:"isBreadcrumb"        gorm:"default:true;comment:是否开启 Breadcrumb 强制经典、横向布局不显示"`
	IsTagsview         bool   `json:"isTagsview"          gorm:"default:true;comment:是否开启 Tagsview"`
	IsTagsviewIcon     bool   `json:"isTagsviewIcon"     gorm:"default:true;comment:是否开启 Tagsview 图标"`
	IsCacheTagsView    bool   `json:"isCacheTagsView"   gorm:"default:true;comment:是否开启 TagsView 缓存"`
	IsSortableTagsView bool   `json:"isSortableTagsView" gorm:"default:true;comment:是否开启 TagsView 拖拽"`
	IsShareTagsView    bool   `json:"isShareTagsView"    gorm:"default:false;comment:是否开启 TagsView 共用"`
	IsFooter           bool   `json:"isFooter"             gorm:"default:true;comment:是否开启 Footer 底部版权信息"`
	IsGrayscale        bool   `json:"isGrayscale"             gorm:"default:false;comment:是否开启灰色模式"`
	IsInvert           bool   `json:"isInvert"             gorm:"default:false;comment:是否开启色弱模式"`
	IsWartermark       bool   `json:"isWartermark"         gorm:"default:false;comment:是否开启水印"`
	WartermarkText     string `json:"wartermarkText"       gorm:"default:AirGo;comment:水印文案"`

	//其它设置
	// Tagsview 风格：可选值"<tags-style-one|tags-style-four|tags-style-five>"，默认 tags-style-five
	// 定义的值与 `/src/layout/navBars/tagsView/tagsView.vue` 中的 class 同名
	TagsStyle          string `json:"tagsStyle"           gorm:"default:tags-style-five;comment:Tagsview 风格：可选值<tags-style-one|tags-style-four|tags-style-five>"`
	Animation          string `json:"animation"            gorm:"default:slide-right;comment:主页面切换动画：可选值<slide-right|slide-left|opacitys>"`
	ColumnsAsideStyle  string `json:"columnsAsideStyle"  gorm:"default:columns-round;comment:分栏高亮风格：可选值<columns-round|columns-card>"`
	ColumnsAsideLayout string `json:"columnsAsideLayout" gorm:"default:columns-horizontal;comment:分栏布局风格：可选值<columns-horizontal|columns-vertical>"`

	//布局切换
	//注意：为了演示，切换布局时，颜色会被还原成默认，代码位置：/@/layout/navBars/breadcrumb/setings.vue
	//中的 `initSetLayoutChange(设置布局切换，重置主题样式)` 方法
	// 布局切换：可选值"<defaults|classic|transverse|columns>"，默认 defaults
	Layout string `json:"layout" gorm:"default:defaults;comment:布局切换：可选值<defaults|classic|transverse|columns>"`

	// 全局网站标题 / 副标题
	// 网站主标题（菜单导航、浏览器当前网页标题）
	GlobalTitle         string `json:"globalTitle"          gorm:"default:AirGo;comment:网站主标题"`
	GlobalViceTitle     string `json:"globalViceTitle"      gorm:"default:Professional and stable!!;comment:网站副标题"`
	GlobalViceTitleMsg  string `json:"globalViceTitleMsg"   gorm:"default:You deserve it;comment:网站副标题"`
	GlobalI18n          string `json:"globalI18n"           gorm:"default:zh-cn;comment:默认初始语言"`
	GlobalComponentSize string `json:"globalComponentSize"  gorm:"default:small;comment:默认全局组件大小，可选值<large|'default'|small>"`
}
