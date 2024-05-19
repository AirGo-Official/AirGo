package constant

const DEFAULT_CLASH_RULE = `rule-providers:
  reject:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/reject.txt"
    path: ./ruleset/reject.yaml
    interval: 86400

  icloud:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/icloud.txt"
    path: ./ruleset/icloud.yaml
    interval: 86400

  apple:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/apple.txt"
    path: ./ruleset/apple.yaml
    interval: 86400

  google:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/google.txt"
    path: ./ruleset/google.yaml
    interval: 86400

  proxy:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/proxy.txt"
    path: ./ruleset/proxy.yaml
    interval: 86400

  direct:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/direct.txt"
    path: ./ruleset/direct.yaml
    interval: 86400

  private:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/private.txt"
    path: ./ruleset/private.yaml
    interval: 86400

  gfw:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/gfw.txt"
    path: ./ruleset/gfw.yaml
    interval: 86400

  tld-not-cn:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/tld-not-cn.txt"
    path: ./ruleset/tld-not-cn.yaml
    interval: 86400

  telegramcidr:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/telegramcidr.txt"
    path: ./ruleset/telegramcidr.yaml
    interval: 86400

  cncidr:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/cncidr.txt"
    path: ./ruleset/cncidr.yaml
    interval: 86400

  lancidr:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/lancidr.txt"
    path: ./ruleset/lancidr.yaml
    interval: 86400

  applications:
    type: http
    behavior: classical
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/applications.txt"
    path: ./ruleset/applications.yaml
    interval: 86400
rules:
  - RULE-SET,applications,DIRECT
  - DOMAIN,clash.razord.top,DIRECT
  - DOMAIN,yacd.haishan.me,DIRECT
  - RULE-SET,private,DIRECT
  - RULE-SET,reject,REJECT
  - RULE-SET,icloud,DIRECT
  - RULE-SET,apple,DIRECT
  - RULE-SET,google,AirGo-PROXY
  - RULE-SET,proxy,AirGo-PROXY
  - RULE-SET,direct,DIRECT
  - RULE-SET,lancidr,DIRECT
  - RULE-SET,cncidr,DIRECT
  - RULE-SET,telegramcidr,AirGo-PROXY
  - GEOIP,LAN,DIRECT
  - GEOIP,CN,DIRECT
  - MATCH,AirGo-PROXY`

const DEFAULT_SURGE_RULE = `PROCESS-NAME,v2ray,DIRECT
PROCESS-NAME,xray,DIRECT
PROCESS-NAME,clash,DIRECT
PROCESS-NAME,naive,DIRECT
PROCESS-NAME,trojan,DIRECT
PROCESS-NAME,trojan-go,DIRECT
PROCESS-NAME,ss-local,DIRECT
PROCESS-NAME,privoxy,DIRECT
PROCESS-NAME,leaf,DIRECT
PROCESS-NAME,Thunder,DIRECT
PROCESS-NAME,DownloadService,DIRECT
PROCESS-NAME,qBittorrent,DIRECT
PROCESS-NAME,Transmission,DIRECT
PROCESS-NAME,fdm,DIRECT
PROCESS-NAME,aria2c,DIRECT
PROCESS-NAME,Folx,DIRECT
PROCESS-NAME,NetTransport,DIRECT
PROCESS-NAME,uTorrent,DIRECT
PROCESS-NAME,WebTorrent,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/private.txt,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/reject.txt,REJECT
RULE-SET,SYSTEM,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/icloud.txt,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/apple.txt,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/google.txt,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/proxy.txt,AirGo-PROXY,force-remote-dns
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/direct.txt,DIRECT
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/telegramcidr.txt,AirGo-PROXY
RULE-SET,https://cdn.jsdelivr.net/gh/Loyalsoldier/surge-rules@release/ruleset/cncidr.txt,DIRECT
RULE-SET,LAN,DIRECT
FINAL,AirGo-PROXY,dns-failed`
