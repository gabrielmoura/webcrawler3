import useVisNetwork from "../../useVisNetwork.ts";
import {useEffect} from "react";
import {Edge, Node, Options} from "vis-network/standalone/esm/vis-network";

interface DataResponse {
    nodes: Node[],
    edges: Edge[],
}

const options: Options = {
    nodes: {
        shape: "dot",
        size: 16,
        scaling: {
            min: 10,
            max: 30,
            // label: {
            //     Min: 8,
            //     Max: 30,
            //     drawThreshold: 12,
            //     maxVisible: 20,
            // },
        },
        font: {
            size: 12,
            face: "Tahoma",
        },
    },
    edges: {
        width: 0.15,
        color: {
            inherit: "from",
            opacity: 0.5,
        },
        // smooth: {
        //     type: "continuous",
        // },
    },
    interaction: {
        tooltipDelay: 200,
        // hideEdgesOnDrag: true,
        // hideEdgesOnZoom: true,
    },
    // layout: {
    //     improvedLayout: false, // Disable the improved layout
    //
    // },
    // physics:false,
    physics: {
        enabled: true,
        barnesHut: {
            gravitationalConstant: -80000,
            centralGravity: 0.3,
            springLength: 95,
            springConstant: 0.04,
            damping: 0.09,
            avoidOverlap: 0
        },
        maxVelocity: 70,
        minVelocity: 0.75,
        solver: "barnesHut",
        stabilization: {
            enabled: true,
            iterations: 1000,
            updateInterval: 1000,
            onlyDynamicEdges: false,
            fit: true
        },
        timestep: 0.5,
        adaptiveTimestep: true
    },
};

export function GraphVis({data}: { data: DataResponse }) {

    const {ref, network} = useVisNetwork({
        options,
        edges: data.edges,
        nodes: data.nodes,
    });

    useEffect(() => {
        if (!network) return;
        network.once("beforeDrawing", () => {
            network.focus(5);
        })

        // network.setSelection({
        //     edges: [1],
        //     nodes: [1]
        // });
    }, [network]);


    return (
        <div style={{height: 700, width: "100%"}} ref={ref}></div>
    );
}