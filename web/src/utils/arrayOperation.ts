/**
 * 判断两数组字符串是否相同（用于按钮权限验证），数组字符串中存在相同时会自动去重（按钮权限标识不会重复）
 * @param news 新数据
 * @param old 源数据
 * @returns 两数组相同返回 `true`，反之则反
 */
export function judementSameArr(newArr: unknown[] | string[], oldArr: string[]): boolean {
    const news = removeDuplicate(newArr);
    const olds = removeDuplicate(oldArr);
    let count = 0;
    const leng = news.length;
    for (let i in olds) {
        for (let j in news) {
            if (olds[i] === news[j]) count++;
        }
    }
    return count === leng ? true : false;
}

/**
 * 判断两个对象是否相同
 * @param a 要比较的对象一
 * @param b 要比较的对象二
 * @returns 相同返回 true，反之则反
 */
export function isObjectValueEqual<T>(a: T, b: T): boolean {
    if (!a || !b) return false;
    let aProps = Object.getOwnPropertyNames(a);
    let bProps = Object.getOwnPropertyNames(b);
    if (aProps.length != bProps.length) return false;
    for (let i = 0; i < aProps.length; i++) {
        let propName = aProps[i];
        // @ts-ignore
        let propA = a[propName];
        // @ts-ignore
        let propB = b[propName];
        if (!b.hasOwnProperty(propName)) return false;
        if (propA instanceof Object) {
            if (!isObjectValueEqual(propA, propB)) return false;
        } else if (propA !== propB) {
            return false;
        }
    }
    return true;
}

/**
 * 数组、数组对象去重
 * @param arr 数组内容
 * @param attr 需要去重的键值（数组对象）
 * @returns
 */
export function removeDuplicate(arr: EmptyArrayType, attr?: string) {
    if (!Object.keys(arr).length) {
        return arr;
    } else {
        if (attr) {
            const obj: EmptyObjectType = {};
            return arr.reduce((cur: EmptyArrayType[], item: EmptyArrayType) => {
                // @ts-ignore
                obj[item[attr]] ? '' : (obj[item[attr]] = true && item[attr] && cur.push(item));
                return cur;
            }, []);
        } else {
            return [...new Set(arr)];
        }
    }
}

//menu数组提取节点
export function arrayExtractionNodes(data: any) {
    let sonArr: any = []
    let sonNewArr: any = []
    let parentArr: any = []
    data.forEach((item: any) => {
        parentArr.push(item.parent_id)
        sonArr.push(item.id)
    });
    //父节点去重
    const newPsrentArr = parentArr.filter((value: any, index: any, array: any) => {
        return array.indexOf(value) === index;
    });
    //子节点去重
    sonArr.forEach((item: any) => {
        if (newPsrentArr.indexOf(item) === -1) {
            sonNewArr.push(item)
        }
    })
    // console.log("newPsrentArr:", newPsrentArr)
    // console.log("sonArr:", sonArr)
    // console.log("sonNewArr:", sonNewArr)
    return sonNewArr
}

// 获取 pinia 中的路由
export const getMenuData = (routes: RouteItems) => {
    const arr: RouteItems = [];
    routes.map((val: RouteItem) => {
        val['title'] = val.meta?.title as string;
        arr.push({...val});
        if (val.children) getMenuData(val.children);
    });
    //console.log("获取 pinia 中的路由 arr:",arr)
    return arr;
};

