import type { ReactNode, FC } from "preact/compat";
import { Navigate, useLocation } from "react-router-dom";
import { signedIn } from "~/contexts/user";
import { checkUser } from "~/helpers";
import { useOnce } from "~/hooks";
import { routes } from "~/router";

const Protected: FC<{ children?: ReactNode }> = ({ children }) => {
    const { pathname } = useLocation()

    useOnce(() => {
        const controller = new AbortController
        checkUser(controller.signal).then((b: boolean) => signedIn.value = b)

        return () => controller.abort()
    })

    if (signedIn.value == undefined) return null;
    if (signedIn.value == false) return <Navigate to={`${routes.users_sign_in}?back=${pathname}`} />;

    return <>{children}</>;
};

export default Protected;
