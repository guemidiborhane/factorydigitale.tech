import {
    ActionFunctionArgs,
    Form,
    LoaderFunctionArgs,
    redirect,
    useActionData,
} from "react-router-dom";
import { APIError, fetchApi } from "~/helpers/http";

import styles from "./index.module.scss";
import Input from "ui/Input";
import { useDirection } from "~/hooks";
import Button from "ui/Button";
import { useT } from "i18n";
import { APISchemas } from "~/api";
import { checkUser } from "~/helpers";
import Toaster from "ui/Toaster";
import { fetchUser } from "~/stores/user";

export async function action({ request }: ActionFunctionArgs) {
    const { signal } = request;

    const body = Object.fromEntries(await request.formData()) as APISchemas['pkg_users.UserRegisterParams']
    const response = await fetchApi("/api/users", {
        method: 'post',
        signal,
        body,
    })

    if (response.ok) {
        console.log('signup, now login')
        const signinResponse = await fetchApi("/api/auth", {
            method: 'post', signal, body
        })

        if (signinResponse.ok) {
            const back = (new URL(request.url)).searchParams.get('back')
            fetchUser()
            return redirect(back || "/")
        }

        return null
    }

    return response.error;
}

export const protect = false
export async function loader({ request }: LoaderFunctionArgs) {
    const response = await checkUser(request.signal)
    const back = (new URL(request.url)).searchParams.get('back')

    if (response) return redirect(back || '/')

    return null
}

export default function LoginPage() {
    const action = useActionData() as APIError
    const { t } = useT()

    return (
        <div dir={useDirection()}>
            <div class={styles.SignUpForm}>
                <Form method="POST" class="box">
                    <Input name="username" label="models.user.username" autofocus />
                    <Input name="password" type="password" label="models.user.password" />
                    <Button type="Success">{t('misc.sign_up')}</Button>
                </Form>
            </div>

            <Toaster />
        </div>
    );
}
