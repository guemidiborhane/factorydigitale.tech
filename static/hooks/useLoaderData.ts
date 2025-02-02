import { useLoaderData as loaderHook } from "react-router-dom";
import type { APIPaths } from "~/api";
import type { FetchResponse } from "~/helpers/http";

export default function useLoaderData<T extends APIPaths>(): FetchResponse<T, { method: 'get' }> {
  return loaderHook() as FetchResponse<T, { method: 'get' }>
}
