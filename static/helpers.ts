import { fetchApi } from "~/helpers/http";
import type { APIEndpoints, APIPaths } from "./api";

type HasDeleteMethod<T> = T extends { method: "delete" } ? true : false
type DeleteAPIPaths = {
    [K in APIPaths]: HasDeleteMethod<APIEndpoints[K]['requests']> extends false ? never : K
}[APIPaths]

export async function destroyAction(url: DeleteAPIPaths, res: Response): Promise<Response | null> {
    const { ok } = await fetchApi(url, { method: 'delete' })

    if (ok) return res

    return null
}

export const csrfToken: string = (() => {
    const cookie = document.cookie.split('; ').find(row => row.startsWith('csrf_='))?.split('=')
    const [, token] = cookie || ["", ""]

    return token
})()

export const checkUser = async (signal?: AbortSignal): Promise<boolean> => {
    const response = await fetchApi('/api/user/check', { signal }, false)

    return response.ok && response.data
}

export const truncate = (str: string, length: number): string => {
    if (str.length <= length) return str;
    return str.slice(0, length) + '...';
}
