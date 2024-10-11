import {createRootRouteWithContext, Outlet,} from '@tanstack/react-router'
import {TanStackRouterDevtools} from '@tanstack/router-devtools'
import '../assets/index.css'
import {QueryClient} from "@tanstack/react-query";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";

export const Route = createRootRouteWithContext<{
    queryClient: QueryClient
}>()({
    component: RootComponent,
})

function RootComponent() {

    return (
        <>

            <Outlet/>

            {import.meta.env.DEV ?
                <>
                    <ReactQueryDevtools buttonPosition="bottom-left"/>
                    <TanStackRouterDevtools position="bottom-right"/>
                </>
                : null}
        </>
    )
}