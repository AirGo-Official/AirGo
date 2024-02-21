package database

const (
	userApiPre  = "/api"
	adminApiPre = "/api/admin"
	// 默认邮件验证码样式
	text1 = `<div >
  <p >欢迎使用，请及时输入验证码</p>
  <span style="font-size:30px">emailcode</span>
</div>`
	// 商品默认描述
	text2 = `<h3 style="color:#00BFFF">究竟什么样的终点，才配得上这一路的颠沛流离---管泽元</h3>
<h3 style="color:#DDA0DD">世界聚焦于你---管泽元</h3>`
	rule1 = `api.vv1234.cn
shua.leyz.top
llss.atewm.cn
siriling.github.io
shidahuilang.github.io
fu-c-k.github.io
db.laomoe.com
loss.98cat.cn
net.ljxnet.cn
`
	rule2 = `(.*\.||)(speedtest)\.(com|cn|net)
(.*\.||)(LibreSpeed)\.(com|cn|net)
(.*\.||)(fast)\.(com|cn|net)
(.*.||)(gov|12377|12315|talk.news.pts.org|cread­ers|zhuich­aguoji|efcc.org|cy­ber­po­lice|abolu­owang|tu­idang|epochtimes|dafa­hao|falundafa|minghui|falu­naz|zhengjian|110.qq|mingjingnews|in­medi­ahk|xin­sheng|banned­book|nt­dtv|falun­gong|12321|se­cretchina|epochweekly|cn.rfi).(cn|com|org|net|club|net|fr|tw|hk)
(..)(dafahao|minghui|falungong|dongtaiwang|epochtimes|ntdtv|falundafa|wujieliulan|zhengjian).(org|com|net) (..)(dafa­hao|minghui|dong­tai­wang|epochtimes|nt­dtv|falundafa|wu­jieli­u­lan|zhengjian).(org|com|net) (.*.)(dafa­hao|minghui|dong­tai­wang|epochtimes|nt­dtv|falundafa|wu­jieli­u­lan|zhengjian).(org|com|net)
(.*\.)(pincong)\.(rocks)
(.*.)(64tianwang|beijingspring|boxun|broadpressinc|chengmingmag|chenpokong|chinaaffairs|chinadigitaltimes|chinesepen|dafahao|dalailamaworld|dalianmeng|dongtaiwang|epochweekly|erabaru|fgmtv|hrichina|huanghuagang|hxwq|jiangweiping|lagranepoca|lantosfoundation|minghui|minzhuzhongguo|ned|ninecommentaries|ogate|renminbao|rfa|secretchina|shenyun|shenyunperformingarts|shenzhoufilm|soundofhope|tiantibooks|tibetpost|truthmoviegroup.wixsite|tuidang|uhrp|uyghuramerican|voachinese|vot|weijingsheng|wujieliulan|xizang-zhiye|zhengjian|zhuichaguoji).(org|com|net)
(.*\.)(gov|12377|12315|talk.news.pts|zhuichaguoji|efcc|cyberpolice|tuidang|falundafa|falunaz|110.qq|mingjingnews|inmediahk|xinsheng|12321|epochweekly|cn.rfi|mingjing|chinaaid|botanwang|xinsheng|rfi|breakgfw|chengmingmag|jinpianwang|xizang-zhiye|breakgfw|qi-gong|voachinese|mhradio|rfa|edoors|edoors|renminbao|soundofhope|zhengjian|dafahao|minghui|dongtaiwang|epochtimes|ntdtv|falundafa|wujieliulan|aboluowang|bannedbook|secretchina|dajiyuan|boxun|chinadigitaltimes|huaglad|dwnews|creaders|oneplusnews|rfa|nextdigital|pincong|gtv|kwok7)\.(cn|com|org|net|club|net|fr|tw|hk|eu|info|me|rocks)
(.*.||)(tuidang.ddns|xinsheng|pincong|pages|newhighlandvision).(cn|com|org|net|club|net|fr|tw|hk|rocks|dev)
`
	defaultHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no"/>
    <meta name="apple-mobile-web-app-capable" content="yes"/>
    <meta name="apple-mobile-web-app-status-bar-style" content="black"/>
    <meta name="apple-touch-fullscreen" content="yes"/>
    <meta name="format-detection" content="email=no"/>
    <meta name="wap-font-scale" content="no"/>
    <meta name="viewport" content="user-scalable=no, width=device-width"/>
    <meta content="telephone=no" name="format-detection"/>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <title>AirGo 稳定、极速的网络</title>
    <link href="https://src.axui.cn/v2.0/dist/css/ax.css" rel="stylesheet" type="text/css">
    <link href="https://src.axui.cn/v2.0/dist/css/ax-response.css" rel="stylesheet" type="text/css">
    <link rel="shortcut icon" href="https://telegraph-image.pages.dev/file/c48a2f45ebf102dd66131.png">

</head>
<body>

<header class="ax-header">
    <div class="ax-row">
        <div class="ax-col"><a href="###" class="ax-logo"><img
                src="https://telegraph-image.pages.dev/file/c48a2f45ebf102dd66131.png"/></a></div>
        <div class="ax-nav">
            <div class="ax-item ax-grade"><a href="/#/home" class="ax-text">首页<i class="ax-iconfont ax-icon-down"></i></a><span
                    class="ax-line"></span></div>
            <div class="ax-item ax-grade"><a href="/#/shop" class="ax-text">商店<i class="ax-iconfont ax-icon-down"></i></a><span
                    class="ax-line"></span></div>
            <div class="ax-item ax-grade"><a href="/#/documents" class="ax-text">文档<i
                    class="ax-iconfont ax-icon-down"></i></a><span class="ax-line"></span></div>
        </div>
    </div>
</header>
<div style=";margin-top: 100px">
    <div class="ax-row ax-split">
        <div class="ax-col ax-col-12">
            <section class="demo-section-border" style="padding: 50px;text-align: center">
                <h1 style="margin-top: 50px;font-size: xxx-large">AirGo 稳定、极速的网络</h1>
                <div style="margin-top:50px;"><span
                        class="ax-glitter ax-primary"></span>专注于全球网络加速的服务商
                </div>
                <div style="margin-top: 50px;">
                    <a href="/#/login" class="ax-btn ax-primary" style="margin-right: 10px"><span
                            class="ax-iconfont ax-icon-triangle-right-f"></span>立即开始</a>
                    <a href="###" class="ax-btn"><span class="ax-iconfont ax-icon-paperplane-f"></span>加入我们</a>
                </div>
            </section>
        </div>
        <div class="ax-col ax-col-12">
            <section class="demo-section-border">
                <svg viewBox="0 0 400 300" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M69.072 230.276h263.351l26.876 43.563H41l28.072-43.563Z" fill="#D2E9FE"/>
                    <path d="M65.013 238.189h175.52" stroke="#398CFF" stroke-width=".686" stroke-miterlimit="10"/>
                    <path d="M230.759 70.036a26.893 26.893 0 0 1 9.935-10.345c4.124-2.352 8.948-3.273 13.313-2.237a7.313 7.313 0 0 1 3.778 1.947c2.242 2.422 1.652 6.691-.13 9.79-1.782 3.097-4.504 5.504-6.576 8.387a24.676 24.676 0 0 0-4.31 10.23c-.675 3.944-.45 8.263-2.737 11.576-1.902 2.748-5.13 4.114-8.223 4.78-3.093.666-6.266.796-9.294 1.832-3.689 1.266-6.927 3.799-10.075 6.286-2.002 1.561-4.004 3.158-6.321 3.974a10.161 10.161 0 0 1-6.201.215 3.953 3.953 0 0 1-2.182-1.341c-.836-1.156-.646-2.873-.08-4.294 1.977-5.005 7.437-7.458 10.555-11.812a14.263 14.263 0 0 0 2.633-7.327c.15-2.297-.286-4.565-.081-6.857a13.095 13.095 0 0 1 5.286-9.269c3.588-2.567 8.548-3.233 10.905-7.072"
                          fill="#D2E9FE" style="mix-blend-mode:multiply"/>
                    <path d="M267.074 97.903c1.086.305 2.247-.125 3.183-.75.936-.626 1.732-1.442 2.673-2.053 2.267-1.466 5.235-1.591 7.797-.73a4.95 4.95 0 0 1 2.002 1.13 3.752 3.752 0 0 1 .501 4.164 5.665 5.665 0 0 1-3.363 2.758 12.1 12.1 0 0 1-4.425.461 39.69 39.69 0 0 1-8.413-1.437c-1.031-.28-2.102-.59-3.158-.42-1.171.19-2.167.936-3.138 1.622-.691.5-1.541 1.001-2.332.68a1.76 1.76 0 0 1-1.001-1.731 3.77 3.77 0 0 1 .786-1.967 6.221 6.221 0 0 1 2.692-2.278c2.002-.785 4.219-.125 6.251.546"
                          fill="#93B9FF" style="mix-blend-mode:multiply"/>
                    <path d="M252.6 119.824a3.473 3.473 0 1 0 .001-6.947 3.473 3.473 0 0 0-.001 6.947ZM356.206 85.956a3.473 3.473 0 1 0 0-6.946 3.473 3.473 0 0 0 0 6.946Z"
                          fill="#2C60FF"/>
                    <path d="M239.723 210.292s-1.707 8.428-3.754 11.731c0 0-2.678 22.322-25.625 26.992-11.561 2.222-33.087-15.015-42.541-32.117 0 0-19.764 22.607-20.96 23.804h-23.919s10.475-18.178 20.56-34.033c7.888-12.408 14.564-22.277 21.816-22.397 7.252-.12 10.135 3.643 12.012 5.125 1.877 1.481 29.734 21.275 29.734 21.275l14.629-18.603 18.048 18.223Z"
                          fill="#1E294A"/>
                    <path d="M146.322 235.126h-19.354M146.818 236.467s14.779-40.694 19.394-39.613c4.614 1.081 26.42 40.429 36.725 40.039"
                          stroke="#fff" stroke-width=".23" stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="M163.619 86.117c-.31-3.334-.906-6.582-.501-9.94a15.015 15.015 0 0 1 3.704-8.318c2.883-3.118 7.477-4.725 11.571-3.594 1.917.53 3.759 1.722 4.55 3.544a6.215 6.215 0 0 1 .33 3.503c-.441 2.473-2.002 4.605-3.764 6.401-1.761 1.797-3.743 3.374-5.355 5.306a18.211 18.211 0 0 0-3.653 7.086c-.501 1.907-.681 3.89-1.257 5.771-.575 1.882-1.621 3.739-3.343 4.695a2.724 2.724 0 0 1-1.957.37 2.467 2.467 0 0 1-1.616-1.857 6.884 6.884 0 0 1 0-2.567c.18-1.247 1.471-8.248 1.291-10.4Z"
                          fill="#398CFF"/>
                    <path d="M225.179 53.555c1.501-3.318 5.385-5.545 8.903-4.63 1.567.41 3.108 1.622 3.163 3.244.04 1.246-.795 2.342-1.681 3.228-.886.886-2.032 1.831-2.778 2.968-1.001 1.501-1.316 3.353-1.857 5.084-.54 1.732-1.456 3.504-3.103 4.27-.565.26-1.386.305-1.671-.256a1.424 1.424 0 0 1-.1-.67v-3.198s-2.348-6.727-.876-10.04Z"
                          fill="#93B9FF"/>
                    <path d="M235.574 43.48c1.246-2.157 3.413-3.954 5.906-4.069a7.15 7.15 0 0 1 2.532.376c1.186.385 2.367 1.06 2.898 2.187.731 1.531.055 3.338-.756 4.83-2.287 4.203-5.705 7.887-6.937 12.511-.525 2.053-.63 4.215-1.026 6.272-.215 1.1-.69 2.357-1.771 2.637a2.186 2.186 0 0 1-1.962-.61 6.167 6.167 0 0 1-1.186-1.767 18.909 18.909 0 0 1-1.732-4.45c-.796-3.673 2.783-15.76 4.034-17.917Z"
                          fill="#AABCFF" style="mix-blend-mode:multiply"/>
                    <path d="M291.523 53.5c5.661 0 10.25-4.59 10.25-10.25S297.184 33 291.523 33s-10.25 4.59-10.25 10.25 4.589 10.25 10.25 10.25Z"
                          fill="#93B9FF"/>
                    <path d="M270.257 56.668a3.874 3.874 0 1 0 0-7.749 3.874 3.874 0 0 0 0 7.749Z" fill="#D2E9FE"
                          style="mix-blend-mode:multiply"/>
                    <path d="M227.581 122.622c1.227 0 2.222-1.454 2.222-3.248 0-1.794-.995-3.248-2.222-3.248-1.227 0-2.222 1.454-2.222 3.248 0 1.794.995 3.248 2.222 3.248Z"
                          fill="#AABCFF" style="mix-blend-mode:multiply"/>
                    <path d="M261.319 54.846a1.371 1.371 0 1 0 0-2.742 1.371 1.371 0 0 0 0 2.742ZM157.743 136.255a5.877 5.877 0 0 1-.395-3.428c.056-.488.266-.946.6-1.306a1.788 1.788 0 0 1 1.757-.326 2.972 2.972 0 0 1 1.437 1.151c.341.523.631 1.078.865 1.657a22.516 22.516 0 0 0 3.444 5.861 17.69 17.69 0 0 0 1.661 1.636c.348.275.662.589.936.936.279.349.413.792.376 1.236a.615.615 0 0 1-.476.571 2.39 2.39 0 0 1-1.912-.285 8.391 8.391 0 0 1-1.501-1.282 32.122 32.122 0 0 0-2.347-2.057c-.414-.332-.846-.64-1.297-.92-.325-.196-.665-.366-.976-.581-1.016-.706-2.172-2.863-2.172-2.863Z"
                          fill="#2C60FF"/>
                    <path d="M157.983 154.998c.396-1.146.766-2.337 1.572-3.253a4.255 4.255 0 0 1 3.188-1.501 2.97 2.97 0 0 1 2.772 2.002c.296 1.001-.055 2.072-.455 3.038-.4.966-.846 1.972-.756 3.003a6.17 6.17 0 0 0 .691 2.002c.344.63.561 1.322.636 2.036a1.93 1.93 0 0 1-.911 1.802 1.873 1.873 0 0 1-1.712-.12 4.63 4.63 0 0 1-1.301-1.201c-.986-1.231-5.676-2.117-3.724-7.808ZM139.575 119.629a4.554 4.554 0 1 0 .001-9.109 4.554 4.554 0 0 0-.001 9.109Z"
                          fill="#93B9FF" style="mix-blend-mode:multiply"/>
                    <path d="M136.102 78.975a1.937 1.937 0 1 0 0-3.874 1.937 1.937 0 0 0 0 3.874Z" fill="#D2E9FE"
                          style="mix-blend-mode:multiply"/>
                    <path d="M183.904 171.354c-3.189-.29-5.556-3.593-7.007-6.151a19.01 19.01 0 0 1-2.528-9.008c0-1.547.221-3.604 1.502-4.675 1.797-1.501 3.503.906 3.954 2.437.92 3.033.035 6.271.56 9.359a13.474 13.474 0 0 0 1.582 4.41c.375.654.805 1.276 1.286 1.857.606.735 1.396 1.346 1.857 2.192l.095-.13a14.175 14.175 0 0 0-1.942-.431c-.11 0-.155.15-.045.165.656.098 1.305.242 1.942.431a.087.087 0 0 0 .048-.004.088.088 0 0 0 .059-.073.087.087 0 0 0-.007-.048 9.728 9.728 0 0 0-1.502-1.757 13.013 13.013 0 0 1-1.656-2.287 13.156 13.156 0 0 1-1.677-5.155c-.34-3.243.781-7.167-1.126-10.09-.646-1.001-1.867-1.897-3.088-1.311-1.301.625-1.777 2.382-1.937 3.688-.405 3.288.596 6.862 2.082 9.775a19.113 19.113 0 0 0 3.218 4.554c1.171 1.221 2.598 2.252 4.325 2.407.11 0 .11-.16 0-.17l.005.015Z"
                          fill="#F0F"/>
                    <path d="M207.737 137.521a1.387 1.387 0 0 1-.941-1.411c.225-3.794 1.947-7.417 4.229-10.48 1.561-2.107 3.503-4.084 6.041-4.81a5.453 5.453 0 0 1 2.717-.13 3.5 3.5 0 0 1 2.202 1.502c.591 1.001.571 2.182.666 3.318.095 1.136.42 2.387 1.396 2.973.661.395 1.472.395 2.223.56.915.205 1.836.751 2.122 1.642a2.691 2.691 0 0 1-.271 2.002 5.676 5.676 0 0 1-3.108 2.682 8.6 8.6 0 0 1-4.139.396c-2.237-.281-4.399-1.206-6.651-1.141a7.703 7.703 0 0 0-4.94 2.002 4.01 4.01 0 0 1-1.141.88.6.6 0 0 1-.405.015Z"
                          fill="#398CFF"/>
                    <path d="M336.172 238.189h-58.057M354.755 268.058l-308.995.02" stroke="#398CFF" stroke-width=".686"
                          stroke-miterlimit="10"/>
                    <path d="M208.662 170.779a9.755 9.755 0 0 1-1.476-4.229c-.04-.421 0-.956.415-1.086a.905.905 0 0 1 .816.285c1.872 1.621 2.112 4.434 3.504 6.506 1.041 1.542 2.672 2.568 4.254 3.549l10.009 6.206-2.132 8.533-16.09-10.4.7-9.364Z"
                          fill="#AABCFF"/>
                    <path d="m146.487 240.702-15.56-.081-1.321 2.988h14.189l2.692-2.907ZM127.944 247.388h14.034l.761-1.581H128.64l-.696 1.581ZM126.353 250.992l14.069-.351.505-1.056h-13.953l-.621 1.407Z"
                          fill="#398CFF"/>
                    <path d="m129.606 243.609-.966 2.197h14.099l1.056-2.197h-14.189ZM126.974 249.585h13.953l1.051-2.197h-14.033l-.971 2.197Z"
                          fill="#F4F4F4"/>
                    <path d="M107.825 261.542c.881 1.406 4.749 1.756 4.749 1.756h35.34v-3.323H107.79a7.672 7.672 0 0 0 .035 1.567Z"
                          fill="#93B9FF"/>
                    <path d="M127.084 249.06s-1.337.37-3.324.941c1.447.655 3.464 1.816 3.849 3.368.616 2.457-3.428 2.457-4.569 2.197-.686-.16-4.615-1.767-7.658-3.033-.871.28-1.726.561-2.542.836-4.054 1.371-4.91 4.639-5.045 6.606h40.124v-13.268c-8.624 5.255-20.835 2.353-20.835 2.353Z"
                          fill="#2C60FF"/>
                    <path d="M127.584 253.369c-.386-1.552-2.403-2.713-3.849-3.368-2.282.66-5.42 1.586-8.378 2.532 3.043 1.266 6.972 2.873 7.657 3.033 1.161.26 5.205.26 4.57-2.197Z"
                          fill="#F4F4F4"/>
                    <path d="m208.753 178.101-10.685-7.558s-1.677-3.728 1.611-4.124c2.172-.255 7.718 2.002 11.441 4.585"
                          fill="#AABCFF"/>
                    <path d="M229.142 172.826a9.927 9.927 0 0 0-.18-2.963 3.343 3.343 0 0 0-1.781-2.267 10.457 10.457 0 0 0-1.177-.336 4.678 4.678 0 0 1-2.107-1.311c-1.191-1.291-1.556-3.143-2.557-4.589a5.927 5.927 0 0 1-.681-1.056c-.345-.806-.135-1.727.09-2.573a77.346 77.346 0 0 1 2.648-8.048l12.697 6.562.696 10.885s-.18 1.692.901 2.092l-8.549 3.604Z"
                          fill="#93B9FF"/>
                    <path d="M222.296 150.309a30.036 30.036 0 0 0 2.928 4.649 10.197 10.197 0 0 0 4.344 3.269c.017-.53.141-1.051.365-1.532a1.72 1.72 0 0 1 1.201-.941c.941-.15 1.767.791 1.862 1.742.095.951-.325 1.862-.75 2.717a1.007 1.007 0 0 0-.151.461c.033.198.09.391.171.575a1.47 1.47 0 0 1-.371 1.287c-.324.333-.666.647-1.026.941-1.601 1.501-1.902 4.059-1.106 6.101l14.304-4.95s-.5-5.921-5.806-6.206c0 0 1.652-8.944-4.044-10.766a14.238 14.238 0 0 0-6.481-.685 8.221 8.221 0 0 0-5.44 3.348v-.01Z"
                          fill="#FF3F08"/>
                    <path d="M236.294 148.762a3.073 3.073 0 1 0 0-6.145 3.073 3.073 0 0 0 0 6.145Z" fill="#FF3F08"/>
                    <path d="M246.229 217.464s-.596 8.924-2.648 12.227c0 0-3.783 21.821-26.73 26.491-23.318 2.037-33.088-15.015-42.542-32.117 0 0-7.752 20.445-8.948 21.641h-23.919s17.207-44.758 22.677-49.142c0 0 2.673-4.955 7.688-5.125 5.015-.17 10.135 3.643 12.011 5.125 1.877 1.481 18.909 17.707 18.909 17.707l25.455-15.015 18.047 18.208Z"
                          fill="#1E294A"/>
                    <path d="M167.843 240.631h-24.279M154.495 240.702s13.608-37.762 18.223-36.681c4.614 1.081 26.42 40.434 36.725 40.039M137.309 235.126s23.082-41.04 27.696-39.969"
                          stroke="#fff" stroke-width=".23" stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="m167.763 175.608 32.627-9.964 16.346 15.029-32.517 12.417-16.456-17.482Z" fill="#398CFF"/>
                    <path d="m172.818 176.524 26.976-8.348 7.572 7.002-26.761 8.944-7.787-7.598Z" fill="#D2E9FE"/>
                    <path d="m190.17 182.73 11.401-3.833 2.953 2.462-11.296 4.304-3.058-2.933Z" fill="#FF3F08"/>
                    <path d="M167.763 175.608v-25.369l32.602-8.453v23.883l-32.602 9.939Z" fill="#398CFF"/>
                    <path d="M169.67 173.041v-20.335l28.643-7.557v19.724l-28.643 8.168Z" fill="#AABCFF"/>
                    <path d="m243.526 170.909-2.447 2.622-1.481-2.452 3.928-.17Z" fill="#2C60FF"/>
                    <path d="M243.046 176.885a33.413 33.413 0 0 0-2.398-5.15l-12.667 2.612c-2.067 1.211-3.708 3.028-5.215 4.895-7.262 9.009-12.172 19.624-16.481 30.344 0 0-7.372-8.468-9.879-8.583-2.508-.115-8.484 5.926-7.633 7.863.851 1.937 16.176 22.091 23.178 20.214 7.001-1.877 14.349-22.437 14.514-23.407 0 0 1.561 7.086 5.505 10.069 9.194-16.08 10.851-34.818 11.076-38.857ZM244.698 182.385l.15 1.001c-.636 13.543-9.895 33.613-9.895 33.613 3.74.444 7.508.596 11.271.455 0 0 2.302-19.294-1.526-35.069Z"
                          fill="#398CFF"/>
                    <path d="M234.958 217.009a6.138 6.138 0 0 1-3.002-1.267s1.436 1.292 3.002 1.267Z" fill="#1E294A"/>
                    <path d="M244.697 182.385a49.086 49.086 0 0 0-1.656-5.505c-.225 4.039-1.882 22.787-11.081 38.857a6.138 6.138 0 0 0 3.003 1.267s9.259-20.07 9.894-33.613l-.16-1.006Z"
                          fill="#1E294A"/>
                    <path d="M225.679 176.049s15.295 17.282 20.735 34.643l-.035-7.292s-8.133-20.79-18.383-29.028c.015-.015-2.107 1.036-2.317 1.677Z"
                          fill="#1E294A"/>
                    <path d="M225.757 160.089c.903-.839 1.203-1.984.671-2.557-.533-.573-1.696-.357-2.599.482-.904.839-1.204 1.984-.671 2.557.532.573 1.696.358 2.599-.482Z"
                          fill="#2C60FF"/>
                    <path d="M226.995 174.993c-.12 0 .901-2.002 1.582-2.463 0 0 11.446-4.894 12.417-4.949 0 0 2.392 2.562 2.562 3.303a69.727 69.727 0 0 1-7.938 3.658c-1.338.48-2.747.727-4.169.731a27.768 27.768 0 0 1-4.454-.28Z"
                          fill="#1E294A"/>
                    <path d="M230.684 159.067s-.2-2.082 1.081-2.452" stroke="#1E294A" stroke-width=".115"
                          stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="M231.185 157.621s.625.5.57 1.081" stroke="#1E294A" stroke-width=".115"
                          stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="m167.763 175.608 32.627-9.909" stroke="#1E294A" stroke-width=".23" stroke-miterlimit="10"/>
                    <path d="M238.287 158.382s5.41-.571 7.292 4.899M235.689 148.697s3.478 1.317 4.73-3.298"
                          stroke="#FF3F08" stroke-width=".17" stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="M180.105 151.265c.445 3.153 2.232 5.976 3.934 8.668.921-3.083 1.882-6.236 3.823-8.803a22.224 22.224 0 0 1 5.005-4.56 62.933 62.933 0 0 1 10.621-6.086c4.589-2.062 9.479-3.603 13.513-6.586a14.133 14.133 0 0 0 4.384-4.91c2.357-4.709.536-10.86-3.553-14.189-4.089-3.328-9.915-4.004-14.98-2.577-5.065 1.426-9.429 4.785-12.918 8.728.026 0-11.941 15.305-9.829 30.315Z"
                          fill="#AABCFF" style="mix-blend-mode:multiply"/>
                    <path d="M188.919 101.902c.501-1.797.501-3.684.621-5.546a27.363 27.363 0 0 1 1.551-8.383c1.031-2.662 2.768-5.12 5.205-6.606a6.137 6.137 0 0 1 3.814-1.001 3.405 3.405 0 0 1 2.948 2.367 5.125 5.125 0 0 1-.375 2.838c-2.818 8.188-6.817 16.761-7.813 25.37-.846 7.287 1.647 15.269.741 22.547-.801 6.401-5.506 10.455-8.674 16.06a3.357 3.357 0 0 1-.856 1.101c-1.076.771-2.692-.095-3.203-1.316-.51-1.221-.29-2.613-.155-3.934a27.439 27.439 0 0 0-1.181-11.301c-1.081-3.328-2.803-6.456-3.503-9.879-1.477-7.007 1.881-14.75 7.927-18.548a7.254 7.254 0 0 0 2.137-2.047 6.063 6.063 0 0 0 .816-1.722Z"
                          fill="#93B9FF" style="mix-blend-mode:multiply"/>
                    <path d="M195.56 163.011a4.822 4.822 0 0 1-1.036 1.402 1.386 1.386 0 0 1-1.621.175 1.538 1.538 0 0 1-.501-.741 3.15 3.15 0 0 1 .531-3.003c.795-.961 2.097-1.356 3.002-2.212 1.207-1.136 1.547-2.888 2.328-4.344a5.086 5.086 0 0 1 2.502-2.433 3.066 3.066 0 0 1 3.313.621c1.001 1.096.806 2.913-.05 4.129-.855 1.216-2.197 2.002-3.503 2.713l-2.502 1.406s-1.992 1.491-2.463 2.287Z"
                          fill="#2C60FF" style="mix-blend-mode:multiply"/>
                    <path d="M206.285 209.576s1.141 7.062 2.618 8.083c1.476 1.021 0-3.073 0-3.073" stroke="#2C60FF"
                          stroke-width=".23" stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="M161.002 245.706h-14.014l-1.366 2.408h14.188l1.192-2.408ZM143.96 251.892h14.034l.76-1.581h-14.098l-.696 1.581ZM142.368 255.496l14.069-.35.505-1.056h-13.953l-.621 1.406Z"
                          fill="#398CFF"/>
                    <path d="m145.622 248.114-.966 2.197h14.098l1.056-2.197h-14.188ZM142.989 254.09h13.954l1.051-2.198H143.96l-.971 2.198Z"
                          fill="#F4F4F4"/>
                    <path d="M119.336 266.046c.881 1.407 4.75 1.757 4.75 1.757h35.339v-3.323h-40.124a7.66 7.66 0 0 0 .035 1.566Z"
                          fill="#93B9FF"/>
                    <path d="M138.594 253.564s-1.336.37-3.323.941c1.447.656 3.464 1.817 3.849 3.368.615 2.458-3.428 2.458-4.57 2.197-.685-.16-4.614-1.766-7.657-3.033-.871.281-1.727.561-2.542.836-4.054 1.372-4.91 4.64-5.045 6.607h40.124v-13.268c-8.624 5.255-20.836 2.352-20.836 2.352Z"
                          fill="#2C60FF"/>
                    <path d="M139.095 257.873c-.385-1.551-2.402-2.712-3.849-3.368-2.282.661-5.42 1.586-8.378 2.532 3.043 1.267 6.972 2.873 7.658 3.033 1.161.261 5.205.261 4.569-2.197Z"
                          fill="#F4F4F4"/>
                    <path d="M176.046 192.235a8.418 8.418 0 0 0 2.317 1.086" stroke="#1E294A" stroke-width=".255"
                          stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="M176.782 177.365a.406.406 0 0 0-.349-.141.398.398 0 0 0-.312.211c-.575 1.086-2.087 4.109-2.026 5.751.07 2.027.795 7.457 1.376 8.543.36.501.575.501.575.501.05.561.147 1.116.291 1.661.145.29 6.441 3.909 8.468 4.925s4.344 3.543 5.43 6.802c0 0 5.005-2.968 6.006-4.705 0 0-3.038-4.199-3.834-5.861-.795-1.661-3.473-6.006-4.199-6.731a70.36 70.36 0 0 0-6.116-4.45c-.58-.29-3.258-2.462-3.763-2.392-.506.07-1.377.22-1.086 1.086.29.866 3.838 4.705 5.285 5.866a5.95 5.95 0 0 1-4.705-1.812 18.624 18.624 0 0 0-.871-2.748s.796-4.344.501-5.215a4.205 4.205 0 0 0-.671-1.291Z"
                          fill="#AABCFF"/>
                    <path d="M176.046 192.235s1.657 1 2.352 1.121" stroke="#1E294A" stroke-width=".23"
                          stroke-miterlimit="10" stroke-linecap="round"/>
                    <path d="M278.721 197.445v21.285h20.77v17.187h29.239v17.687h-68.222v-56.159h18.213Z"
                          fill="#AABCFF"/>
                    <path d="M255.443 253.969h-18.818l27.647-99.661h18.818l-27.647 99.661Z" fill="#93B9FF"
                          style="mix-blend-mode:multiply"/>
                </svg>
            </section>
        </div>
    </div>
</div>

<div style="text-align: center;margin-top: 50px">
    <h1>产品优势</h1>
    <h4>全球网络数据中心，支撑全时段高速连接</h4>
</div>

<div style="margin-left: 15%;margin-top: 50px;width: 1000px;" class="container">
<img src="https://telegraph-image.pages.dev/file/117a2efa6f3f38516a63b.png">
<div class="box" style="left: 160px;top: 230px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 140px;top: 350px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 300px;top: 450px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 500px;top: 400px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 700px;top: 380px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 800px;top: 390px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 800px;top: 530px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 550px;top: 260px;"><span class="ax-glitter ax-primary"></span></div>
<div class="box" style="left: 480px;top: 300px;"><span class="ax-glitter ax-primary"></span></div>
</div>

<div class="ax-footer" style="margin: 50px">
    <div class="ax-row ax-split">
        <div class="ax-col ax-copyright">AirGo版权所有　备案号:粤ICP备1234567号　<a href="###"><img
                src="http://www.beian.gov.cn/img/new/gongan.png"/>粤公网安备SB110号</a></div>
        <div class="ax-operate">
            <a href="###" class="ax-iconfont ax-icon-qq-f"></a>
            <a href="###" class="ax-iconfont ax-icon-wechat-f"></a>
            <a href="###" class="ax-iconfont ax-icon-facebook-f"></a>
            <a href="###" class="ax-iconfont ax-icon-twitter-f"></a>
            <a href="###" class="ax-iconfont ax-icon-paperplane-f"></a>
        </div>
    </div>
</div>


<script src="https://src.axui.cn/v2.0/dist/js/ax.min.js" type="text/javascript" charset="utf-8"></script>

<script type='text/javascript'>

</script>
</body>
<style>
    .container {
        position: relative;
    }
    .box {
        /*border: 2px solid red;*/
        position: absolute;
        /*z-index: -1;*/
    }
</style>
</html>
`
	defaultDialog = `# 欢迎使用！最新活动如下：
1. 免费注册体验！注册送流量
2. 邀请返利

**更多套餐请查看商店**`
)
