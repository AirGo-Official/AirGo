/**
 * 后端时间字符串时区转换，"2023-05-29T17:28:47.50276+08:00" ---> "2023/05/29 17:28:47"
 * @param strTime
 * @constructor
 */
export function DateStrToTime(strTime: string) {
    const date = new Date(strTime);
    return date.toLocaleString();
}

/**
 * 后端时间字符串转换 "2023-05-29T17:28:47.50276+08:00" ---> "2023-05-29"
 * @param strTime
 * @constructor
 */
export function DateStrHandler(strTime: string) {
    const date = new Date(strTime);
    return formatDate(date,'YYYY-mm-dd')
}

/**
 * 某个时间与当前时间的相差的毫秒数
 * @param startTime
 * @param endTime
 * @constructor
 */
export function getTimeDifference(startTime: string) {
    const start = new Date(startTime)
    const end = new Date()
    return  end.getTime() - start.getTime();//时间差的毫秒数
}


/**
 * 获取某个月的时间范围
 * @param n,n=0当月，n=1上月
 * @constructor
 */
export function GetDurationMonth(n: number) {
    let currentDate = new Date();
    let Y = currentDate.getFullYear();
    let M = currentDate.getMonth();
    // let MonthDayNum = new Date(currentY,currentM,0).getDate();  //计算当月的天数
    let startDate = new Date(Y, M - n, 1);
    let endDate = new Date(Y, M + 1 - n, 0, 23, 59, 59);
    // console.log("开始：", startDate)
    // console.log("结束：", endDate)
    let start = formatDate(startDate, "YYYY-mm-dd HH:MM:SS")
    let end = formatDate(endDate, "YYYY-mm-dd HH:MM:SS")
    return [start,end]
}

/**
 * 获取某个月第一天到当前天的时间范围，m为第几月，范围1～12
 * @param m
 * @constructor
 */
export function GetDurationMonthFromFirstToCurrent(m: number) {
    let currentDate = new Date();
    let Y = currentDate.getFullYear();
    let d = currentDate.getDate()
    let startDate = new Date(Y, m+1, 1);
    let endDate = new Date(Y, m + 1, d, 23, 59, 59);
    // console.log("开始：", startDate)
    // console.log("结束：", endDate)
    let start = formatDate(startDate, "YYYY-mm-dd HH:MM:SS")
    let end = formatDate(endDate, "YYYY-mm-dd HH:MM:SS")
    return [start,end]
}

/**
 * 获取当天时间范围
 * @constructor
 */
export function GetDurationDay() {
    let currentDate = new Date();
    let Y = currentDate.getFullYear();
    let M = currentDate.getMonth();
    let D = currentDate.getDate()
    // let MonthDayNum = new Date(currentY,currentM,0).getDate();  //计算当月的天数
    let startDate = new Date(Y, M, D);
    let endDate = new Date(Y, M, D, 23, 59, 59);
    // console.log("开始：", startDate)
    // console.log("结束：", endDate)
    let start = formatDate(startDate, "YYYY-mm-dd HH:MM:SS")
    let end = formatDate(endDate, "YYYY-mm-dd HH:MM:SS")
    return [start,end]
}

/**
 * 获取n天前到今天的时间间隔
 * @param n
 * @constructor
 */
export function GetDurationWithNDay(n:number) {
    let currentDate = new Date();
    let Y = currentDate.getFullYear();
    let M = currentDate.getMonth();
    let D = currentDate.getDate()
    let startDate = new Date(Y, M, D-n+1);
    let endDate = new Date(Y, M, D, 23, 59, 59);
    // console.log("开始：", startDate)
    // console.log("结束：", endDate)
    let start = formatDate(startDate, "YYYY-mm-dd HH:MM:SS")
    let end = formatDate(endDate, "YYYY-mm-dd HH:MM:SS")
    return [start,end]
}

/**
 * 得到标准时区的时间的函数，参数i为时区值数字，比如北京为东八区则输进8,西5输入-5
 * @param i
 * @constructor
 */
export function GetLocalTime(i: number) {
    if (typeof i !== 'number') return;
    var d = new Date();
    //得到1970年一月一日到现在的秒数
    var len = d.getTime();
    //本地时间与GMT时间的时间偏移差(注意：GMT这是UTC的民间名称。GMT=UTC）
    var offset = d.getTimezoneOffset() * 60000;
    //得到现在的格林尼治时间
    var utcTime = len + offset;
    return new Date(utcTime + 3600000 * i);
}

/**
 * 时间日期转换
 * @param date 当前时间，new Date() 格式
 * @param format 需要转换的时间格式字符串
 * @description format 字符串随意，如 `YYYY-mm、YYYY-mm-dd`
 * @description format 季度："YYYY-mm-dd HH:MM:SS QQQQ"
 * @description format 星期："YYYY-mm-dd HH:MM:SS WWW"
 * @description format 几周："YYYY-mm-dd HH:MM:SS ZZZ"
 * @description format 季度 + 星期 + 几周："YYYY-mm-dd HH:MM:SS WWW QQQQ ZZZ"
 * @returns 返回拼接后的时间字符串
 */
export function formatDate(date: Date, format: string): string {
    let we = date.getDay(); // 星期
    let z = getWeek(date); // 周
    let qut = Math.floor((date.getMonth() + 3) / 3).toString(); // 季度
    const opt: { [key: string]: string } = {
        'Y+': date.getFullYear().toString(), // 年
        'm+': (date.getMonth() + 1).toString(), // 月(月份从0开始，要+1)
        'd+': date.getDate().toString(), // 日
        'H+': date.getHours().toString(), // 时
        'M+': date.getMinutes().toString(), // 分
        'S+': date.getSeconds().toString(), // 秒
        'q+': qut, // 季度
    };
    // 中文数字 (星期)
    const week: { [key: string]: string } = {
        '0': '日',
        '1': '一',
        '2': '二',
        '3': '三',
        '4': '四',
        '5': '五',
        '6': '六',
    };
    // 中文数字（季度）
    const quarter: { [key: string]: string } = {
        '1': '一',
        '2': '二',
        '3': '三',
        '4': '四',
    };
    if (/(W+)/.test(format))
        format = format.replace(RegExp.$1, RegExp.$1.length > 1 ? (RegExp.$1.length > 2 ? '星期' + week[we] : '周' + week[we]) : week[we]);
    if (/(Q+)/.test(format)) format = format.replace(RegExp.$1, RegExp.$1.length == 4 ? '第' + quarter[qut] + '季度' : quarter[qut]);
    if (/(Z+)/.test(format)) format = format.replace(RegExp.$1, RegExp.$1.length == 3 ? '第' + z + '周' : z + '');
    for (let k in opt) {
        let r = new RegExp('(' + k + ')').exec(format);
        // 若输入的长度不为1，则前面补零
        if (r) format = format.replace(r[1], RegExp.$1.length == 1 ? opt[k] : opt[k].padStart(RegExp.$1.length, '0'));
    }
    return format;
}

/**
 * 获取当前日期是第几周
 * @param dateTime 当前传入的日期值
 * @returns 返回第几周数字值
 */
export function getWeek(dateTime: Date): number {
    let temptTime = new Date(dateTime.getTime());
    // 周几
    let weekday = temptTime.getDay() || 7;
    // 周1+5天=周六
    temptTime.setDate(temptTime.getDate() - weekday + 1 + 5);
    let firstDay = new Date(temptTime.getFullYear(), 0, 1);
    let dayOfWeek = firstDay.getDay();
    let spendDay = 1;
    if (dayOfWeek != 0) spendDay = 7 - dayOfWeek + 1;
    firstDay = new Date(temptTime.getFullYear(), 0, 1 + spendDay);
    let d = Math.ceil((temptTime.valueOf() - firstDay.valueOf()) / 86400000);
    let result = Math.ceil(d / 7);
    return result;
}

/**
 * 将时间转换为 `几秒前`、`几分钟前`、`几小时前`、`几天前`
 * @param param 当前时间，new Date() 格式或者字符串时间格式
 * @param format 需要转换的时间格式字符串
 * @description param 10秒：  10 * 1000
 * @description param 1分：   60 * 1000
 * @description param 1小时： 60 * 60 * 1000
 * @description param 24小时：60 * 60 * 24 * 1000
 * @description param 3天：   60 * 60* 24 * 1000 * 3
 * @returns 返回拼接后的时间字符串
 */
export function formatPast(param: string | Date, format: string = 'YYYY-mm-dd'): string {
    // 传入格式处理、存储转换值
    let t: any, s: number;
    // 获取js 时间戳
    let time: number = new Date().getTime();
    // 是否是对象
    typeof param === 'string' || 'object' ? (t = new Date(param).getTime()) : (t = param);
    // 当前时间戳 - 传入时间戳
    time = Number.parseInt(`${time - t}`);
    if (time < 10000) {
        // 10秒内
        return '刚刚';
    } else if (time < 60000 && time >= 10000) {
        // 超过10秒少于1分钟内
        s = Math.floor(time / 1000);
        return `${s}秒前`;
    } else if (time < 3600000 && time >= 60000) {
        // 超过1分钟少于1小时
        s = Math.floor(time / 60000);
        return `${s}分钟前`;
    } else if (time < 86400000 && time >= 3600000) {
        // 超过1小时少于24小时
        s = Math.floor(time / 3600000);
        return `${s}小时前`;
    } else if (time < 259200000 && time >= 86400000) {
        // 超过1天少于3天内
        s = Math.floor(time / 86400000);
        return `${s}天前`;
    } else {
        // 超过3天
        let date = typeof param === 'string' || 'object' ? new Date(param) : param;
        return formatDate(date, format);
    }
}

/**
 * 时间问候语
 * @param param 当前时间，new Date() 格式
 * @description param 调用 `formatAxis(new Date())` 输出 `上午好`
 * @returns 返回拼接后的时间字符串
 */
export function formatAxis(param: Date): string {
    let hour: number = new Date(param).getHours();
    if (hour < 6) return '凌晨好';
    else if (hour < 9) return '早上好';
    else if (hour < 12) return '上午好';
    else if (hour < 14) return '中午好';
    else if (hour < 17) return '下午好';
    else if (hour < 19) return '傍晚好';
    else if (hour < 22) return '晚上好';
    else return '夜里好';
}

// 获取指定格式时间
export function getFormatDate(fmt: string, ts: string = '') {
    const date = ts ? new Date(ts) : new Date()
    let o: Record<string, any> = {
        'M+': date.getMonth() + 1,
        'd+': date.getDate(),
        'H+': date.getHours(),
        'm+': date.getMinutes(),
        's+': date.getSeconds(),
        'q+': Math.floor((date.getMonth() + 3) / 3),
        'S': date.getMilliseconds()
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
    for (let k in o) {
        let item = o[k];
        if (new RegExp('(' + k + ')').test(fmt))
            fmt = fmt.replace(RegExp.$1, RegExp.$1.length == 1 ? item : ('00' + item).substr(('' + item).length))
    }
    return fmt
}


