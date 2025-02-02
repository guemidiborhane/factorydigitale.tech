import { deepMap, onMount } from "nanostores";
import { APIResponse } from "~/api";
import { fetchApi } from "~/helpers/http";

type Response = APIResponse<'/api/user', 'get'>
export type UserStore = Response['user'] & {
  permissions: Response['permissions'],
  tracks: Response['tracks']
}
export const $user = deepMap<UserStore>()

export async function fetchUser() {
  const response = await fetchApi("/api/user")

  if (!response.ok) return null

  const { user, permissions, tracks } = response.data

  $user.set({
    ...user,
    permissions,
    tracks
  })
}

onMount($user, () => {
  fetchUser()
});
