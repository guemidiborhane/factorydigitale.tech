import {
    type ActionFunctionArgs,
    type LoaderFunctionArgs,
    RouterProvider,
    createBrowserRouter,
    Outlet,
} from "react-router-dom";
import Protected from "ui/Protected";
import { type FC, StrictMode } from "preact/compat";
import { UserProvider } from "./contexts/user";

type LazyComponent = {
    default: FC;
    protect: boolean;
    action?: (args: ActionFunctionArgs) => Response | null;
    loader: (args: LoaderFunctionArgs) => Response;
};

type PreservedRoutes = Record<string, () => Promise<LazyComponent>>;

// @ts-ignore
const ROUTES = import.meta.glob("/pkg/**/pages/**/[a-z[]*.ts(x)?");
// @ts-ignore
const PRESERVED = import.meta.glob("/pkg/**/pages/(_app|404).tsx");

const lazyLoader = (promise: () => Promise<LazyComponent>, protect = true) => async () => {
    const { default: Component, protect: p, ...props } = await promise();
    protect = typeof p == 'boolean' ? p : protect
    const ProtectedComponent = () => (
        <Protected>
            <Component />
        </Protected>
    );

    return {
        Component: protect ? ProtectedComponent : Component,
        ...props
    };
};

const normalize = (str: string): string => {
    return str
        .replace(/\/pkg\/|\/static\//, "")
        .replace(/\/pages\//, "/")
        .replace(/\/index/, "")
        .replace(/\.tsx?$/, "")
        .replace(/\[(.+?)\]/g, ":$1")
        .replace(/\[\.{3}(.+?)\]/g, "*$1");
};

const preserved: PreservedRoutes = Object.keys(PRESERVED).reduce(
    (preserved, file) => ({
        ...preserved,
        [normalize(file)]: lazyLoader(PRESERVED[file], false),
    }),
    {},
);

const getPreserved = async (name: string) => {
    const layout = preserved[name] && (await preserved[name]());
    return layout || { Component: Outlet };
};

const children = Object.keys(ROUTES)
    .filter((file) => !file.endsWith(".module.scss.d.ts"))
    .map((file) => {
        const [pkg, ...p] = normalize(file).split("/");
        const path = p.join("/");
        const lazy = lazyLoader(ROUTES[file]);

        return { path, lazy, pkg, file };
    });

const childrenByPkg: Record<string, any[]> = children.reduce(
    (group, child) => {
        const { pkg } = child;
        group[pkg] = group[pkg] ?? [];
        group[pkg].push(child);

        return group;
    },
    {} as Record<string, any>,
);

const r = [
    {
        path: "/",
        async lazy() {
            const { default: Component } = await import("./Home");
            return {
                Component: () => (
                    <Protected>
                        <Component />
                    </Protected>
                ),
            };
        },
    },

    ...Object.keys(childrenByPkg).map((pkg) => ({
        path: `/${pkg}`,
        async lazy() {
            return await getPreserved(`${pkg}/_app`);
        },
        children: childrenByPkg[pkg].map((route) => {
            const index = route.path == ''
            return { ...route, index }
        }),
    })),

    {
        path: "*",
        async lazy() {
            return await getPreserved("404");
        },
    },
];

export const routes = children.reduce(
    (output, child) => {
        const { pkg, path, file } = child
        const join = [pkg, path].filter((a) => a != "").join('_')
        output[join] = `/${join.replace(/_/, '/')}`

        if (process.env.NODE_ENV === "development") {
            output[`${join}_file`] = file
        }

        return output
    },
    {} as Record<string, any>
)

if (process.env.NODE_ENV === "development") {
    console.groupCollapsed('Routes')
    console.table(Object.keys(routes).filter(r => !r.endsWith('_file')).map(name => ({
        name,
        route: routes[name],
        file: routes[`${name}_file`]
    })))
    console.groupEnd()
}

const router = createBrowserRouter(r);
export function Router() {
    return (
        <StrictMode>
            <UserProvider>
                <RouterProvider router={router} />
            </UserProvider>
        </StrictMode>
    );
}
