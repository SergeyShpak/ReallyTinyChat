export interface IRequest<T> {
    method: string; path: string; code: number,
    params?: any, body?: any,
    request?: () => Uint8Array,
    response?: (responseBinary: Uint8Array) => T,
    before?: (r: XMLHttpRequest, body: Uint8Array) => void,
    host?: string
}

export class Client {

    public debug: boolean
    private host: string

    constructor(host: string) {
        this.host = host
    }

    public doRequest<T>(r: IRequest<T>): Promise<T> {
        const host = r.host != null ? r.host : this.host
        return new Promise((resolve, reject) => {
            const xmlhttp = new XMLHttpRequest()
            const uri = host + this.uri_query(r.path, r.params)
            xmlhttp.open(r.method, uri, true)
            xmlhttp.responseType = 'arraybuffer'
            xmlhttp.setRequestHeader("Cache-Control", "no-cache")
            xmlhttp.onreadystatechange = (e) => {
                if (xmlhttp.readyState !== 4 || xmlhttp.status === 0) {
                    return
                }
                if (this.debug) {
                    console.debug("response(" + r.method + "," + host + r.path + "): ", xmlhttp.status)
                }
                if (xmlhttp.status !== r.code) {
                  return reject({
                      code: xmlhttp.status
                  })
                }
                if (r.response == null) {
                    if (xmlhttp.response === null || xmlhttp.response.length === 0) {
                        console.debug("WARNING: response is not used", xmlhttp.response)
                    }
                    return resolve()
                }
                let resp
                try {
                    resp = JSON.parse(xmlhttp.response)
                } catch (e) {
                    return reject(e)
                }
                resolve(resp)
            }
            xmlhttp.onerror = (e) => {
                reject(e)
            }
            if (this.debug) {
                console.debug("request(" + r.method + "," + host + r.path + ")")
            }
            console.log("Body: ", JSON.stringify(r.body))
            xmlhttp.send(JSON.stringify(r.body))
        })
    }

    private uri_query(url: string, params?: any): string {
        if (params == null) {
            return url
        }
        const uriParams = Object.keys(params).map(key => {
            const x = key + "=";
            return ((params[key].constructor === Array) ?
                x + params[key].map(encodeURIComponent).join("&" + x) :
                x + encodeURIComponent(params[key]))
        }).join("&")
        if (uriParams.length === 0) {
            return url
        }
        return url + "?" + uriParams
    }
}
