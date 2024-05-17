// import NodeRSA from "node-rsa"
//
// // RSA加密
// function RSA(pkey:string, data:string) {
//     // let buff = Buffer.from(data);
//     let key = new NodeRSA(pkey);
//     key.setOptions({encryptionScheme: 'pkcs1'});
//     // let encrypted = key.encrypt(buff, 'base64', 'base64');
//     return key.encrypt(data, 'base64');
// }
// Vite不使用node-rsa模块,使用 jsencryp


import NodeRSA from 'jsencrypt';
// RSA加密
function RSA(pkey:string, data:string) {
    let key = new NodeRSA();
    key.setPublicKey(pkey)
    return key.encrypt(data).toString();
}

// 电信手机号加密处理
export  function TelecomMobileHandler(phoneNum: string)
{
    let result = ''
    let l = phoneNum.length
    let ArrPhone = phoneNum.toString().split('')
    for (let i = 0; i < l; i++) {
        result = `${result}` + String.fromCharCode(ArrPhone[i].charCodeAt(0) + 2 & 65535)
    }
    return result
}

//电信 RSA加密
export function TelecomRSAEncrypt(str: string){
    const key = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDBkLT15ThVgz6/NOl6s8GNPofdWzWbCkWnkaAm7O2LjkM1H7dMvzkiqdxU02jamGRHLX/ZNMCXHnPcW/sDhiFCBN18qFvy8g6VYb9QtroI09e176s+ZCtiv7hbin2cCTj99iUpnEloZm19lwHyo69u5UMiPMpq0/XKBO8lYhN/gwIDAQAB";
    return RSA(key,str)
}


// 随机字符串,不包含大写字母
export function randomStringNoUpper(length:number) {
    // const chars = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const chars = '123456789abcdefghijklmnopqrstuvwxyz';
    let result = '';
    for (let i = length; i > 0; --i) result += chars[Math.floor(Math.random() * chars.length)];
    return result;
}
// 随机字符串,包含大写字母
export function randomStringWithUpper(length:number) {
    const chars = '123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    let result = '';
    for (let i = length; i > 0; --i) result += chars[Math.floor(Math.random() * chars.length)];
    return result;
}

//随机邀请码，以用户id开头
export function randomInvitation(userID:number,length:number) {
    return userID+randomStringWithUpper(length)
}
