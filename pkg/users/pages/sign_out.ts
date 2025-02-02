import { redirectDocument } from "react-router-dom"
import { destroyAction } from "~/helpers"

export async function action() {
  return await destroyAction(
    '/api/auth',
    redirectDocument('/user/sign_in')
  )
}
