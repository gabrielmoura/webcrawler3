import {createFileRoute} from '@tanstack/react-router'
import axios from "axios";

import {useQuery} from "@tanstack/react-query";
import {Graph2} from "../../components/netwrokd3/Graph.tsx";


export const Route = createFileRoute('/graph/')({
    component: ShowGraph,
})


function ShowGraph() {

    const {isLoading, data} = useQuery({
        queryKey: ['graph'],
        queryFn: async () => {
            return axios.get('http://localhost:9090/api/v1/manager/graph').then((response) => {
                return response.data
            })
        },
    })

    if (isLoading) return <div>Loading...</div>

    if (data) return <Graph2 data={data}/>
    return <div>Failed to load data</div>

}








