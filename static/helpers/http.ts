import toast from "react-hot-toast";
import type { APIPaths, APIRequests, APIResponse, APISchemas } from "~/api";
import { csrfToken } from "~/helpers";

export type APIError = APISchemas["github_com_guemidiborhane_factorydigitale_tech_internal_errors.HttpError"]

export type FetchResponse<Path extends APIPaths, Options extends APIRequests<Path>> =
  | {
    ok: true
    data: APIResponse<Path, Options['method']>,
  }
  | {
    ok: false
    error: APIError,
  }

export async function fetchApi<
  Path extends APIPaths,
  Options extends APIRequests<Path>
>(url: Path, params?: Options & { signal?: AbortSignal }, notify: boolean = true): Promise<FetchResponse<Path, Options>> {
  const fetchOptions: RequestInit = {
    method: params?.method ?? 'get',
    credentials: 'include',
    headers: {
      "Content-type": "application/json; charset=UTF-8",
      "Accept": "application/json",
      "X-CSRF-Token": csrfToken
    }
  }
  const options = (params ?? {}) as Options

  const body = 'body' in options ? options['body'] : null
  if (body) {
    (fetchOptions.headers as Record<string, string>)['Content-Type'] = 'application/json'
    fetchOptions.body = JSON.stringify(body)
  }

  let urlPath: string = url

  if ('urlParams' in options) {
    for (const [name, value] of Object.entries(options.urlParams)) {
      urlPath = urlPath.replace(`{${name}}`, value.toString())
    }
  }


  if ('query' in options) {
    let searchParams: URLSearchParams = new URLSearchParams()
    for (const [name, value] of Object.entries(options.query || {})) {
      searchParams.set(name, value.toString())
    }

    urlPath = `${url.toString()}?${searchParams.toString()}`
  }


  const shouldNotify = notify && !url.startsWith('/api/user')
  let toastId: undefined | string = undefined

  if (shouldNotify) {
    toastId = toast.loading(url)
  }
  try {
    const response = await fetch(urlPath, fetchOptions)
    if (!response.headers.get('content-type')?.includes('application/json')) {
      if (shouldNotify) toast.error(url, { id: toastId })
      console.error("Response is not JSON")
    }

    const data = await response.json()

    if (response.ok) {
      if (shouldNotify) toast.success(url, { id: toastId })
      return {
        ok: true,
        data
      }
    } else {
      if (shouldNotify) toast.error(url, { id: toastId })
      return {
        ok: false,
        error: data
      }
    }

  } catch (error) {
    if (shouldNotify) toast.error(url, { id: toastId })

    return {
      ok: false,
      error: error as APIError,
    }
  }

}
