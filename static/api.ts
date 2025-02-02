export type APIRequestBodies = {
  "pkg_permissions.PermissionsParams": APISchemas["pkg_permissions.PermissionsParams"]
}

export type APISchemas = {
  "github_com_guemidiborhane_factorydigitale_tech_internal_errors.HttpError": {
    code: string
    message: string | { [key: string]: string }
    status: number
  }
  "github_com_guemidiborhane_factorydigitale_tech_pkg_movies_models.Favourite": {
    id: number
    movie_id: number
    user_id: number
  }
  "github_com_guemidiborhane_factorydigitale_tech_pkg_movies_models.Movie": {
    genres: Array<string>
    id: number
    in_favourites: boolean
    overview: string
    poster: string
    release_date: number
    title: string
  }
  "github_com_guemidiborhane_factorydigitale_tech_pkg_users_auth.UserResponse": {
    id: number
    role: string
    username: string
  }
  "github_com_guemidiborhane_factorydigitale_tech_pkg_users_auth.UserTracks": {
    ips?: Array<Array<string>>
    login_count?: number
  }
  "pkg_movies.FavouriteRequestParams": { movie_id?: number }
  "pkg_permissions.PermissionsParams": { [key: string]: Array<string> }
  "pkg_permissions.PermissionsResponse": {
    permissions: { [key: string]: Array<string> }
    roles_permissions: APISchemas["pkg_permissions.RolesPermissionsMap"]
  }
  "pkg_permissions.RolesPermissionsMap": {
    [key: string]: { [key: string]: Array<string> }
  }
  "pkg_users.UserJSONResponse": {
    permissions: Array<string>
    tracks: APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_users_auth.UserTracks"]
    user: APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_users_auth.UserResponse"]
  }
  "pkg_users.UserRegisterParams": { password: string; username: string }
  "pkg_users_auth.LoginParams": { password: string; username: string }
  "pkg_users_auth.UserResponse": { id: number; role: string; username: string }
}

export type APIEndpoints = {
  "/api/auth": {
    responses: { delete: null; post: APISchemas["pkg_users_auth.UserResponse"] }
    requests:
      | { method: "delete" }
      | { method: "post"; body: APISchemas["pkg_users_auth.LoginParams"] }
  }
  "/api/movies": {
    responses: {
      get: Array<
        APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_movies_models.Movie"]
      >
    }
    requests: { method?: "get"; query?: { offset?: number } }
  }
  "/api/movies/favourite": {
    responses: {
      post: APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_movies_models.Favourite"]
    }
    requests: {
      method: "post"
      body: APISchemas["pkg_movies.FavouriteRequestParams"]
    }
  }
  "/api/movies/favourites": {
    responses: {
      get: Array<
        APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_movies_models.Movie"]
      >
    }
    requests: { method?: "get" }
  }
  "/api/permissions": {
    responses: { get: APISchemas["pkg_permissions.PermissionsResponse"] }
    requests: { method?: "get" }
  }
  "/api/permissions/{role}": {
    responses: { post: Array<string>; put: Array<string> }
    requests:
      | {
          method: "post"
          urlParams: { role: string }
          body: APISchemas["pkg_permissions.PermissionsParams"]
        }
      | {
          method: "put"
          urlParams: { role: string }
          body: APISchemas["pkg_permissions.PermissionsParams"]
        }
  }
  "/api/user": {
    responses: { get: APISchemas["pkg_users.UserJSONResponse"] }
    requests: { method?: "get" }
  }
  "/api/user/check": {
    responses: { get: boolean }
    requests: { method?: "get" }
  }
  "/api/users": {
    responses: {
      post: APISchemas["github_com_guemidiborhane_factorydigitale_tech_pkg_users_auth.UserResponse"]
    }
    requests: {
      method: "post"
      body: APISchemas["pkg_users.UserRegisterParams"]
    }
  }
}

export type APIPaths = keyof APIEndpoints

export type APIRequests<T extends APIPaths> = APIEndpoints[T]["requests"]

export type APIMethods<T extends APIPaths> = NonNullable<
  APIRequests<T>["method"]
>

export type APIRequest<T extends APIPaths, M extends APIMethods<T>> = Omit<
  {
    [MM in APIMethods<T>]: APIRequests<T> & { method: MM }
  }[M],
  "method"
> & { method?: M }

type DefaultToGet<T extends string | undefined> = T extends string ? T : "get"

export type APIResponse<
  T extends APIPaths,
  M extends string | undefined
> = DefaultToGet<M> extends keyof APIEndpoints[T]["responses"]
  ? APIEndpoints[T]["responses"][DefaultToGet<M>]
  : never
