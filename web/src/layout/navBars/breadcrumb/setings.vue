<template>
  <div class="layout-breadcrumb-seting">
    <el-drawer :title="$t('message.layout.configTitle')" v-model="state.isDrawer" direction="rtl" destroy-on-close size="260px">
      <el-scrollbar class="layout-breadcrumb-seting-bar">
        <el-divider content-position="left">{{$t('message.layout.websiteSettings')}}</el-divider>
        <div class="layout-breadcrumb-seting-bar-flex mt14">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{$t('message.layout.globalTitle')}}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-input v-model="getThemeConfig.globalTitle" size="default" style="width: 150px"></el-input>
          </div>
        </div>

        <div class="layout-breadcrumb-seting-bar-flex mt14">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{$t('message.layout.globalViceTitle')}}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-input v-model="getThemeConfig.globalViceTitle" size="default" style="width: 150px"></el-input>
          </div>
        </div>

        <div class="layout-breadcrumb-seting-bar-flex mt14">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{$t('message.layout.globalViceTitleMsg')}}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-input v-model="getThemeConfig.globalViceTitleMsg" size="default" style="width: 150px"></el-input>
          </div>
        </div>

        <div class="layout-breadcrumb-seting-bar-flex mt14">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{$t('message.layout.logo_link')}}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-input v-model="getThemeConfig.logo_link" size="default" style="width: 150px"></el-input>
          </div>
        </div>
        <!-- 全局主题 -->
        <el-divider content-position="left">{{ $t('message.layout.oneTitle') }}</el-divider>
        <div class="layout-breadcrumb-seting-bar-flex">
          <div class="layout-breadcrumb-seting-bar-flex-label">primary</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-color-picker v-model="getThemeConfig.primary" size="default" @change="onColorPickerChange"> </el-color-picker>
          </div>
        </div>
        <!--        深色模式-->
<!--        <div class="layout-breadcrumb-seting-bar-flex mt15">-->
<!--          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsDark') }}</div>-->
<!--          <div class="layout-breadcrumb-seting-bar-flex-value">-->
<!--            <el-switch v-model="getThemeConfig.isIsDark" size="small" @change="onAddDarkChange"></el-switch>-->
<!--          </div>-->
<!--        </div>-->

        <!-- 顶栏设置 -->
        <el-divider content-position="left">{{ $t('message.layout.twoTopTitle') }}</el-divider>
        <!--        顶栏背景-->
        <div class="layout-breadcrumb-seting-bar-flex">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoTopBar') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-color-picker v-model="getThemeConfig.topBar" size="default" @change="onBgColorPickerChange('topBar')"> </el-color-picker>
          </div>
        </div>
        <!--        顶栏默认字体颜色-->
        <div class="layout-breadcrumb-seting-bar-flex">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoTopBarColor') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-color-picker v-model="getThemeConfig.topBarColor" size="default" @change="onBgColorPickerChange('topBarColor')"> </el-color-picker>
          </div>
        </div>
        <!--        顶栏背景渐变-->
        <div class="layout-breadcrumb-seting-bar-flex mt10">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoIsTopBarColorGradual') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isTopBarColorGradual" size="small" @change="onTopBarGradualChange"></el-switch>
          </div>
        </div>

        <!-- 菜单设置 -->
        <el-divider content-position="left">{{ $t('message.layout.twoMenuTitle') }}</el-divider>
        <!--        菜单背景-->
        <div class="layout-breadcrumb-seting-bar-flex">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoMenuBar') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-color-picker v-model="getThemeConfig.menuBar" size="default" @change="onBgColorPickerChange('menuBar')"> </el-color-picker>
          </div>
        </div>
        <!--        菜单默认字体颜色-->
        <div class="layout-breadcrumb-seting-bar-flex">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoMenuBarColor') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-color-picker v-model="getThemeConfig.menuBarColor" size="default" @change="onBgColorPickerChange('menuBarColor')"> </el-color-picker>
          </div>
        </div>
        <!--        菜单高亮背景色-->
        <div class="layout-breadcrumb-seting-bar-flex">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoMenuBarActiveColor') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-color-picker
              v-model="getThemeConfig.menuBarActiveColor"
              size="default"
              show-alpha
              @change="onBgColorPickerChange('menuBarActiveColor')"
            />
          </div>
        </div>
        <!--        菜单背景渐变-->
        <div class="layout-breadcrumb-seting-bar-flex mt14">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoIsMenuBarColorGradual') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isMenuBarColorGradual" size="small" @change="onMenuBarGradualChange"></el-switch>
          </div>
        </div>

        <!-- 分栏设置 -->
        <div>
          <el-divider content-position="left" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">{{$t('message.layout.twoColumnsTitle') }}</el-divider>
          <!--        分栏菜单背景-->
          <div class="layout-breadcrumb-seting-bar-flex" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoColumnsMenuBar') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-color-picker
                v-model="getThemeConfig.columnsMenuBar"
                size="default"
                @change="onBgColorPickerChange('columnsMenuBar')"
                :disabled="getThemeConfig.layout !== 'columns'"
              >
              </el-color-picker>
            </div>
          </div>
          <!--        分栏菜单默认字体颜色-->
          <div class="layout-breadcrumb-seting-bar-flex" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoColumnsMenuBarColor') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-color-picker
                v-model="getThemeConfig.columnsMenuBarColor"
                size="default"
                @change="onBgColorPickerChange('columnsMenuBarColor')"
                :disabled="getThemeConfig.layout !== 'columns'"
              >
              </el-color-picker>
            </div>
          </div>
          <!--        分栏菜单背景渐变-->
          <div class="layout-breadcrumb-seting-bar-flex mt14" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoIsColumnsMenuBarColorGradual') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch
                v-model="getThemeConfig.isColumnsMenuBarColorGradual"
                size="small"
                @change="onColumnsMenuBarGradualChange"
                :disabled="getThemeConfig.layout !== 'columns'"
              ></el-switch>
            </div>
          </div>
          <!--        分栏菜单鼠标悬停预加载-->
          <div class="layout-breadcrumb-seting-bar-flex mt14" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.twoIsColumnsMenuHoverPreload') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch
                v-model="getThemeConfig.isColumnsMenuHoverPreload"
                size="small"
                @change="onColumnsMenuHoverPreloadChange"
                :disabled="getThemeConfig.layout !== 'columns'"
              ></el-switch>
            </div>
          </div>
        </div>


        <!-- 界面设置 -->
        <div>
          <el-divider content-position="left">{{ $t('message.layout.threeTitle') }}</el-divider>
          <!--        菜单水平折叠-->
          <div class="layout-breadcrumb-seting-bar-flex" :style="{ opacity: getThemeConfig.layout === 'transverse' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.threeIsCollapse') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch
                v-model="getThemeConfig.isCollapse"
                :disabled="getThemeConfig.layout === 'transverse'"
                size="small"
                @change="onThemeConfigChange"
              ></el-switch>
            </div>
          </div>
          <!--        菜单手风琴-->
          <div class="layout-breadcrumb-seting-bar-flex mt15" :style="{ opacity: getThemeConfig.layout === 'transverse' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.threeIsUniqueOpened') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch
                v-model="getThemeConfig.isUniqueOpened"
                :disabled="getThemeConfig.layout === 'transverse'"
                size="small"
                @change="setLocalThemeConfig"
              ></el-switch>
            </div>
          </div>
          <!--        固定 Header-->
          <div class="layout-breadcrumb-seting-bar-flex mt15">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.threeIsFixedHeader') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch v-model="getThemeConfig.isFixedHeader" size="small" @change="onIsFixedHeaderChange"></el-switch>
            </div>
          </div>
          <!--        经典布局分割菜单-->
          <div class="layout-breadcrumb-seting-bar-flex mt15" :style="{ opacity: getThemeConfig.layout !== 'classic' ? 0.5 : 1 }">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.threeIsClassicSplitMenu') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch
                v-model="getThemeConfig.isClassicSplitMenu"
                :disabled="getThemeConfig.layout !== 'classic'"
                size="small"
                @change="onClassicSplitMenuChange"
              >
              </el-switch>
            </div>
          </div>
          <!--        开启锁屏-->
          <div class="layout-breadcrumb-seting-bar-flex mt15">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.threeIsLockScreen') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-switch v-model="getThemeConfig.isLockScreen" size="small" @change="setLocalThemeConfig"></el-switch>
            </div>
          </div>
          <!--        自动锁屏(s/秒)-->
          <div class="layout-breadcrumb-seting-bar-flex mt11">
            <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.threeLockScreenTime') }}</div>
            <div class="layout-breadcrumb-seting-bar-flex-value">
              <el-input-number
                v-model="getThemeConfig.lockScreenTime"
                controls-position="right"
                :min="1"
                :max="9999"
                @change="setLocalThemeConfig"
                size="default"
                style="width: 90px"
              >
              </el-input-number>
            </div>
          </div>

        </div>


        <!-- 界面显示 -->
        <el-divider content-position="left">{{ $t('message.layout.fourTitle') }}</el-divider>
        <!--        侧边栏 Logo-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsShowLogo') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isShowLogo" size="small" @change="onIsShowLogoChange"></el-switch>
          </div>
        </div>

        <div
          class="layout-breadcrumb-seting-bar-flex mt15"
          :style="{ opacity: getThemeConfig.layout === 'classic' || getThemeConfig.layout === 'transverse' ? 0.5 : 1 }"
        >
          <!--        开启 Breadcrumb-->
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsBreadcrumb') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch
              v-model="getThemeConfig.isBreadcrumb"
              :disabled="getThemeConfig.layout === 'classic' || getThemeConfig.layout === 'transverse'"
              size="small"
              @change="onIsBreadcrumbChange"
            ></el-switch>
          </div>
        </div>
        <!--        开启 Tagsview-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsTagsview') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isTagsview" size="small" @change="setLocalThemeConfig"></el-switch>
          </div>
        </div>
        <!--        开启 Tagsview 图标-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsTagsviewIcon') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isTagsviewIcon" size="small" @change="setLocalThemeConfig"></el-switch>
          </div>
        </div>
        <!--        开启 TagsView 缓存-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsCacheTagsView') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isCacheTagsView" size="small" @change="setLocalThemeConfig"></el-switch>
          </div>
        </div>
        <!--        开启 TagsView 拖拽-->
        <div class="layout-breadcrumb-seting-bar-flex mt15" :style="{ opacity: state.isMobile ? 0.5 : 1 }">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsSortableTagsView') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch
              v-model="getThemeConfig.isSortableTagsView"
              :disabled="state.isMobile ? true : false"
              size="small"
              @change="onSortableTagsViewChange"
            ></el-switch>
          </div>
        </div>
        <!--        开启 TagsView 共用-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsShareTagsView') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isShareTagsView" size="small" @change="onShareTagsViewChange"></el-switch>
          </div>
        </div>
        <!--        开启 Footer-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsFooter') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isFooter" size="small" @change="setLocalThemeConfig"></el-switch>
          </div>
        </div>
        <!--        灰色模式-->
<!--        <div class="layout-breadcrumb-seting-bar-flex mt15">-->
<!--          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsGrayscale') }}</div>-->
<!--          <div class="layout-breadcrumb-seting-bar-flex-value">-->
<!--            <el-switch v-model="getThemeConfig.isGrayscale" size="small" @change="onAddFilterChange('grayscale')"></el-switch>-->
<!--          </div>-->
<!--        </div>-->
<!--        &lt;!&ndash;        色弱模式&ndash;&gt;-->
<!--        <div class="layout-breadcrumb-seting-bar-flex mt15">-->
<!--          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsInvert') }}</div>-->
<!--          <div class="layout-breadcrumb-seting-bar-flex-value">-->
<!--            <el-switch v-model="getThemeConfig.isInvert" size="small" @change="onAddFilterChange('invert')"></el-switch>-->
<!--          </div>-->
<!--        </div>-->
        <!--        开启水印-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourIsWartermark') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-switch v-model="getThemeConfig.isWartermark" size="small" @change="onWartermarkChange"></el-switch>
          </div>
        </div>
        <!--        水印文案-->
        <div class="layout-breadcrumb-seting-bar-flex mt14">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fourWartermarkText') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-input v-model="getThemeConfig.wartermarkText" size="default" style="width: 90px" @input="onWartermarkTextInput"></el-input>
          </div>
        </div>

        <!-- 其它设置 -->
        <el-divider content-position="left">{{ $t('message.layout.fiveTitle') }}</el-divider>
        <!--        Tagsview 风格-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fiveTagsStyle') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-select v-model="getThemeConfig.tagsStyle" placeholder="请选择" size="default" style="width: 90px" @change="setLocalThemeConfig">
              <el-option label="风格1" value="tags-style-one"></el-option>
              <el-option label="风格4" value="tags-style-four"></el-option>
              <el-option label="风格5" value="tags-style-five"></el-option>
            </el-select>
          </div>
        </div>
        <!--        主页面切换动画-->
        <div class="layout-breadcrumb-seting-bar-flex mt15">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fiveAnimation') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-select v-model="getThemeConfig.animation" placeholder="请选择" size="default" style="width: 90px" @change="setLocalThemeConfig">
              <el-option label="slide-right" value="slide-right"></el-option>
              <el-option label="slide-left" value="slide-left"></el-option>
              <el-option label="opacitys" value="opacitys"></el-option>
            </el-select>
          </div>
        </div>
        <!--        分栏高亮风格-->
        <div class="layout-breadcrumb-seting-bar-flex mt15" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fiveColumnsAsideStyle') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-select
              v-model="getThemeConfig.columnsAsideStyle"
              placeholder="请选择"
              size="default"
              style="width: 90px"
              :disabled="getThemeConfig.layout !== 'columns' ? true : false"
              @change="setLocalThemeConfig"
            >
              <el-option label="圆角" value="columns-round"></el-option>
              <el-option label="卡片" value="columns-card"></el-option>
            </el-select>
          </div>
        </div>
        <!--        分栏布局风格-->
        <div class="layout-breadcrumb-seting-bar-flex mt15 mb27" :style="{ opacity: getThemeConfig.layout !== 'columns' ? 0.5 : 1 }">
          <div class="layout-breadcrumb-seting-bar-flex-label">{{ $t('message.layout.fiveColumnsAsideLayout') }}</div>
          <div class="layout-breadcrumb-seting-bar-flex-value">
            <el-select
              v-model="getThemeConfig.columnsAsideLayout"
              placeholder="请选择"
              size="default"
              style="width: 90px"
              :disabled="getThemeConfig.layout !== 'columns' ? true : false"
              @change="setLocalThemeConfig"
            >
              <el-option label="水平" value="columns-horizontal"></el-option>
              <el-option label="垂直" value="columns-vertical"></el-option>
            </el-select>
          </div>
        </div>

        <!-- 布局切换 -->
        <el-divider content-position="left">{{ $t('message.layout.sixTitle') }}</el-divider>
        <div class="layout-drawer-content-flex">
          <!-- defaults 布局 -->
          <div class="layout-drawer-content-item" @click="onSetLayout('defaults')">
            <section class="el-container el-circular" :class="{ 'drawer-layout-active': getThemeConfig.layout === 'defaults' }">
              <aside class="el-aside" style="width: 20px"></aside>
              <section class="el-container is-vertical">
                <header class="el-header" style="height: 10px"></header>
                <main class="el-main"></main>
              </section>
            </section>
            <div class="layout-tips-warp" :class="{ 'layout-tips-warp-active': getThemeConfig.layout === 'defaults' }">
              <div class="layout-tips-box">
                <p class="layout-tips-txt">{{ $t('message.layout.sixDefaults') }}</p>
              </div>
            </div>
          </div>
          <!-- classic 布局 -->
          <div class="layout-drawer-content-item" @click="onSetLayout('classic')">
            <section class="el-container is-vertical el-circular" :class="{ 'drawer-layout-active': getThemeConfig.layout === 'classic' }">
              <header class="el-header" style="height: 10px"></header>
              <section class="el-container">
                <aside class="el-aside" style="width: 20px"></aside>
                <section class="el-container is-vertical">
                  <main class="el-main"></main>
                </section>
              </section>
            </section>
            <div class="layout-tips-warp" :class="{ 'layout-tips-warp-active': getThemeConfig.layout === 'classic' }">
              <div class="layout-tips-box">
                <p class="layout-tips-txt">{{ $t('message.layout.sixClassic') }}</p>
              </div>
            </div>
          </div>
          <!-- transverse 布局 -->
          <div class="layout-drawer-content-item" @click="onSetLayout('transverse')">
            <section class="el-container is-vertical el-circular" :class="{ 'drawer-layout-active': getThemeConfig.layout === 'transverse' }">
              <header class="el-header" style="height: 10px"></header>
              <section class="el-container">
                <section class="el-container is-vertical">
                  <main class="el-main"></main>
                </section>
              </section>
            </section>
            <div class="layout-tips-warp" :class="{ 'layout-tips-warp-active': getThemeConfig.layout === 'transverse' }">
              <div class="layout-tips-box">
                <p class="layout-tips-txt">{{ $t('message.layout.sixTransverse') }}</p>
              </div>
            </div>
          </div>
          <!-- columns 布局 -->
          <div class="layout-drawer-content-item" @click="onSetLayout('columns')">
            <section class="el-container el-circular" :class="{ 'drawer-layout-active': getThemeConfig.layout === 'columns' }">
              <aside class="el-aside-dark" style="width: 10px"></aside>
              <aside class="el-aside" style="width: 20px"></aside>
              <section class="el-container is-vertical">
                <header class="el-header" style="height: 10px"></header>
                <main class="el-main"></main>
              </section>
            </section>
            <div class="layout-tips-warp" :class="{ 'layout-tips-warp-active': getThemeConfig.layout === 'columns' }">
              <div class="layout-tips-box">
                <p class="layout-tips-txt">{{ $t('message.layout.sixColumns') }}</p>
              </div>
            </div>
          </div>
        </div>


        <div>
          <el-button size="default" type="info" @click="state.isDrawer = false">{{$t('message.common.button_cancel')}}</el-button>
          <el-button size="default" type="primary" @click="onSubmit">{{$t('message.common.button_confirm')}}</el-button>
        </div>
      </el-scrollbar>
    </el-drawer>
  </div>
</template>

<script setup lang="ts" name="layoutBreadcrumbSeting">
import {nextTick, onUnmounted, onMounted, computed, reactive, onBeforeMount} from 'vue';
import {ElMessage} from 'element-plus';
import {storeToRefs} from 'pinia';
import {useThemeConfig} from '/@/stores/themeConfig';
import {useChangeColor} from '/@/utils/theme';
import {verifyAndSpace} from '/@/utils/toolsValidate';
import {Local,Session} from '/@/utils/storage';
import Watermark from '/@/utils/watermark';
import commonFunction from '/@/utils/commonFunction';
import other from '/@/utils/other';
import mittBus from '/@/utils/mitt';
import { useI18n } from 'vue-i18n';


// 定义变量内容
const { locale } = useI18n();
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const {copyText} = commonFunction();
const {getLightColor, getDarkColor} = useChangeColor();
const state = reactive({
  isMobile: false,
  isDrawer: false,  //是否打开主题setting侧窗
});

// 获取布局配置信息
const getThemeConfig = computed(() => {
  return themeConfig.value;
});
// 1、全局主题
const onColorPickerChange = () => {
  if (!getThemeConfig.value.primary) return ElMessage.warning('全局主题 primary 颜色值不能为空');
  // 颜色加深
  document.documentElement.style.setProperty('--el-color-primary-dark-2', `${getDarkColor(getThemeConfig.value.primary, 0.1)}`);
  document.documentElement.style.setProperty('--el-color-primary', getThemeConfig.value.primary);
  // 颜色变浅
  for (let i = 1; i <= 9; i++) {
    document.documentElement.style.setProperty(`--el-color-primary-light-${i}`, `${getLightColor(getThemeConfig.value.primary, i / 10)}`);
  }
  setDispatchThemeConfig();
};
// 2、菜单 / 顶栏
const onBgColorPickerChange = (bg: string) => {
  document.documentElement.style.setProperty(`--next-bg-${bg}`, themeConfig.value[bg]);
  if (bg === 'menuBar') {
    document.documentElement.style.setProperty(`--next-bg-menuBar-light-1`, getLightColor(getThemeConfig.value.menuBar, 0.05));
  }
  onTopBarGradualChange();
  onMenuBarGradualChange();
  onColumnsMenuBarGradualChange();
  setDispatchThemeConfig();
};
// 2、菜单 / 顶栏 --> 顶栏背景渐变
const onTopBarGradualChange = () => {
  setGraduaFun('.layout-navbars-breadcrumb-index', getThemeConfig.value.isTopBarColorGradual, getThemeConfig.value.topBar);
};
// 2、菜单 / 顶栏 --> 菜单背景渐变
const onMenuBarGradualChange = () => {
  setGraduaFun('.layout-container .el-aside', getThemeConfig.value.isMenuBarColorGradual, getThemeConfig.value.menuBar);
};
// 2、菜单 / 顶栏 --> 分栏菜单背景渐变
const onColumnsMenuBarGradualChange = () => {
  setGraduaFun('.layout-container .layout-columns-aside', getThemeConfig.value.isColumnsMenuBarColorGradual, getThemeConfig.value.columnsMenuBar);
};
// 2、菜单 / 顶栏 --> 背景渐变函数
const setGraduaFun = (el: string, bool: boolean, color: string) => {
  nextTick(() => {
    setTimeout(() => {
      let els = document.querySelector(el);
      if (!els) return false;
      document.documentElement.style.setProperty('--el-menu-bg-color', document.documentElement.style.getPropertyValue('--next-bg-menuBar'));
      if (bool) els.setAttribute('style', `background:linear-gradient(to bottom , ${color}, ${getLightColor(color, 0.5)})`);
      else els.setAttribute('style', ``);
      setLocalThemeConfig();
    }, 300);
  });
};
// 2、分栏设置 ->
const onColumnsMenuHoverPreloadChange = () => {
  setLocalThemeConfig();
};
// 3、界面设置 --> 菜单水平折叠
const onThemeConfigChange = () => {
  setDispatchThemeConfig();
};
// 3、界面设置 --> 固定 Header
const onIsFixedHeaderChange = () => {
  getThemeConfig.value.isFixedHeaderChange = getThemeConfig.value.isFixedHeader ? false : true;
  setLocalThemeConfig();
};
// 3、界面设置 --> 经典布局分割菜单
const onClassicSplitMenuChange = () => {
  getThemeConfig.value.isBreadcrumb = false;
  setLocalThemeConfig();
  mittBus.emit('getBreadcrumbIndexSetFilterRoutes');
};
// 4、界面显示 --> 侧边栏 Logo
const onIsShowLogoChange = () => {
  getThemeConfig.value.isShowLogoChange = getThemeConfig.value.isShowLogo ? false : true;
  setLocalThemeConfig();
};
// 4、界面显示 --> 面包屑 Breadcrumb
const onIsBreadcrumbChange = () => {
  if (getThemeConfig.value.layout === 'classic') {
    getThemeConfig.value.isClassicSplitMenu = false;
  }
  setLocalThemeConfig();
};
// 4、界面显示 --> 开启 TagsView 拖拽
const onSortableTagsViewChange = () => {
  mittBus.emit('openOrCloseSortable');
  setLocalThemeConfig();
};
// 4、界面显示 --> 开启 TagsView 共用
const onShareTagsViewChange = () => {
  mittBus.emit('openShareTagsView');
  setLocalThemeConfig();
};
// 4、界面显示 --> 灰色模式/色弱模式
const onAddFilterChange = (attr: string) => {
  if (attr === 'grayscale') {
    if (getThemeConfig.value.isGrayscale) getThemeConfig.value.isInvert = false;
  } else {
    if (getThemeConfig.value.isInvert) getThemeConfig.value.isGrayscale = false;
  }
  const cssAttr =
      attr === 'grayscale' ? `grayscale(${getThemeConfig.value.isGrayscale ? 1 : 0})` : `invert(${getThemeConfig.value.isInvert ? '80%' : '0%'})`;
  const appEle = document.body;
  appEle.setAttribute('style', `filter: ${cssAttr}`);
  setLocalThemeConfig();
};
// 4、界面显示 --> 深色模式
const onAddDarkChange = () => {
  const body = document.documentElement as HTMLElement;
  if (getThemeConfig.value.isIsDark) body.setAttribute('data-theme', 'dark');
  else body.setAttribute('data-theme', '');
};
// 4、界面显示 --> 开启水印
const onWartermarkChange = () => {
  getThemeConfig.value.isWartermark ? Watermark.set(getThemeConfig.value.wartermarkText) : Watermark.del();
  setLocalThemeConfig();


};
// 4、界面显示 --> 水印文案
const onWartermarkTextInput = (val: string) => {
  getThemeConfig.value.wartermarkText = verifyAndSpace(val);
  if (getThemeConfig.value.wartermarkText === '') return false;
  if (getThemeConfig.value.isWartermark) Watermark.set(getThemeConfig.value.wartermarkText);
  setLocalThemeConfig();
};
// 5、布局切换
const onSetLayout = (layout: string) => {
  Session.set('oldLayout', layout);
  if (getThemeConfig.value.layout === layout) return false;
  if (layout === 'transverse') getThemeConfig.value.isCollapse = false;
  getThemeConfig.value.layout = layout;
  // state.isDrawer = false;
  initLayoutChangeFun();
};
// 设置布局切换函数
const initLayoutChangeFun = () => {
  onBgColorPickerChange('menuBar');
  onBgColorPickerChange('menuBarColor');
  onBgColorPickerChange('menuBarActiveColor');
  onBgColorPickerChange('topBar');
  onBgColorPickerChange('topBarColor');
  onBgColorPickerChange('columnsMenuBar');
  onBgColorPickerChange('columnsMenuBarColor');
};
//保存
const onSubmit = () => {
  //请求服务器修改theme
  storesThemeConfig.updateThemeConfig(getThemeConfig.value)
  state.isDrawer = false;
};
// 布局配置弹窗打开
const openDrawer = () => {
  state.isDrawer = true;
};
// 触发 store 布局配置更新
const setDispatchThemeConfig = () => {
  setLocalThemeConfig();
  setLocalThemeConfigStyle();
};
// 存储布局配置
const setLocalThemeConfig = () => {
  Session.remove('themeConfig');
  Session.set('themeConfig', getThemeConfig.value);
};
// 存储布局配置全局主题样式（html根标签）
const setLocalThemeConfigStyle = () => {
  Session.set('themeConfigStyle', document.documentElement.style.cssText);
};
// 初始化菜单样式等
const initSetStyle = () => {
  // 2、菜单 / 顶栏 --> 顶栏背景渐变
  onTopBarGradualChange();
  // 2、菜单 / 顶栏 --> 菜单背景渐变
  onMenuBarGradualChange();
  // 2、菜单 / 顶栏 --> 分栏菜单背景渐变
  onColumnsMenuBarGradualChange();
};

onMounted(() => {
  nextTick(() => {
    // 判断当前布局是否不相同，不相同则初始化当前布局的样式，防止监听窗口大小改变时，布局配置logo、菜单背景等部分布局失效问题
    if (!Session.get('frequency')) initLayoutChangeFun();
    Session.set('frequency', 1);
    // 监听窗口大小改变，非默认布局，设置成默认布局（适配移动端）
    mittBus.on('layoutMobileResize', (res: LayoutMobileResize) => {
      getThemeConfig.value.layout = res.layout;
      state.isDrawer = false;
      initLayoutChangeFun();
      state.isMobile = other.isMobile();
    });
    setTimeout(() => {
      // 默认样式
      onColorPickerChange();
      // 灰色模式
      if (getThemeConfig.value.isGrayscale) onAddFilterChange('grayscale');
      // 色弱模式
      if (getThemeConfig.value.isInvert) onAddFilterChange('invert');
      // 深色模式
      if (getThemeConfig.value.isIsDark) onAddDarkChange();
      // 开启水印
      onWartermarkChange();
      // 语言国际化
      if (Session.get('themeConfig')) locale.value = Session.get('themeConfig').globalI18n;
      // 初始化菜单样式等
      initSetStyle();
    }, 100);
  });
});
onUnmounted(() => {
  mittBus.off('layoutMobileResize', () => {
  });
});

// 暴露变量
defineExpose({
  openDrawer,
});
</script>

<style scoped lang="scss">
.layout-breadcrumb-seting-bar {
  height: calc(100vh - 50px);
  padding: 0 15px;

  :deep(.el-scrollbar__view) {
    overflow-x: hidden !important;
  }

  .layout-breadcrumb-seting-bar-flex {
    display: flex;
    align-items: center;
    margin-bottom: 5px;

    &-label {
      flex: 1;
      color: var(--el-text-color-primary);
    }
  }

  .layout-drawer-content-flex {
    overflow: hidden;
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    margin: 0 -5px;

    .layout-drawer-content-item {
      width: 50%;
      height: 70px;
      cursor: pointer;
      border: 1px solid transparent;
      position: relative;
      padding: 5px;

      .el-container {
        height: 100%;

        .el-aside-dark {
          background-color: var(--next-color-seting-header);
        }

        .el-aside {
          background-color: var(--next-color-seting-aside);
        }

        .el-header {
          background-color: var(--next-color-seting-header);
        }

        .el-main {
          background-color: var(--next-color-seting-main);
        }
      }

      .el-circular {
        border-radius: 2px;
        overflow: hidden;
        border: 1px solid transparent;
        transition: all 0.3s ease-in-out;
      }

      .drawer-layout-active {
        border: 1px solid;
        border-color: var(--el-color-primary);
      }

      .layout-tips-warp,
      .layout-tips-warp-active {
        transition: all 0.3s ease-in-out;
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        border: 1px solid;
        border-color: var(--el-color-primary-light-5);
        border-radius: 100%;
        padding: 4px;

        .layout-tips-box {
          transition: inherit;
          width: 30px;
          height: 30px;
          z-index: 9;
          border: 1px solid;
          border-color: var(--el-color-primary-light-5);
          border-radius: 100%;

          .layout-tips-txt {
            transition: inherit;
            position: relative;
            top: 5px;
            font-size: 12px;
            line-height: 1;
            letter-spacing: 2px;
            white-space: nowrap;
            color: var(--el-color-primary-light-5);
            text-align: center;
            transform: rotate(30deg);
            left: -1px;
            background-color: var(--next-color-seting-main);
            width: 32px;
            height: 17px;
            line-height: 17px;
          }
        }
      }

      .layout-tips-warp-active {
        border: 1px solid;
        border-color: var(--el-color-primary);

        .layout-tips-box {
          border: 1px solid;
          border-color: var(--el-color-primary);

          .layout-tips-txt {
            color: var(--el-color-primary) !important;
            background-color: var(--next-color-seting-main) !important;
          }
        }
      }

      &:hover {
        .el-circular {
          transition: all 0.3s ease-in-out;
          border: 1px solid;
          border-color: var(--el-color-primary);
        }

        .layout-tips-warp {
          transition: all 0.3s ease-in-out;
          border-color: var(--el-color-primary);

          .layout-tips-box {
            transition: inherit;
            border-color: var(--el-color-primary);

            .layout-tips-txt {
              transition: inherit;
              color: var(--el-color-primary) !important;
              background-color: var(--next-color-seting-main) !important;
            }
          }
        }
      }
    }
  }

  .copy-config {
    margin: 10px 0;

    .copy-config-btn {
      width: 100%;
      margin-top: 15px;
    }

    .copy-config-btn-reset {
      width: 100%;
      margin: 10px 0 0;
    }
  }
}
</style>
