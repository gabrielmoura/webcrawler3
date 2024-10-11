import {useEffect, useRef} from "react";

import * as d3 from 'd3';

export interface GraphD3Props {
    data: DataResponse;
    linkStroke?: string;
    linkStrokeOpacity?: number;
    linkStrokeWidth?: number;
    linkStrokeLinecap?: string;
    linkStrength?: number;
    colors?: readonly string[];
    nodeStroke?: string; // Node stroke color
    nodeStrokeWidth?: number; // Node stroke width
    nodeStrokeOpacity?: number; // Node stroke opacity
    nodeRadius?: number; // Node radius
}

export interface DataResponse {
    nodes: Node[],
    edges: Edge[],
}

interface Node {
    id: string | number;
    label: string;
    group?: number;
}

interface Edge {
    from: string | number;
    to: string | number;
}

interface NodeProps extends Node {
    x: number;
    y: number;
}

interface LinkProps {
    source: string | number;
    target: string | number;
}


// Zoom, drag e fixar
export function Graph({
                          data,
                          linkStroke = "#999",
                          linkStrokeOpacity = 0.6,
                          linkStrokeWidth = 1.5,
                          linkStrokeLinecap = "round",
                          linkStrength = 1,
                          colors = d3.schemeTableau10,
                          nodeStroke = "#fff",
                          nodeStrokeWidth = 1.5,
                          nodeStrokeOpacity = 1,
                          nodeRadius = 5,
                      }: GraphD3Props) {
    const svgRef = useRef<SVGSVGElement | null>(null);

    useEffect(() => {
        const drawGraph = (data: DataResponse) => {
            const links: LinkProps[] = data.edges.map(edge => ({
                source: edge.from,
                target: edge.to,
            }));
            const nodes: NodeProps[] = data.nodes.map(node => ({
                id: node.id,
                x: 0,
                y: 0,
                label: node.label,
                group: node.group || 0,
            }));

            const width = 800;
            const height = 600;

            const svg = d3.select(svgRef.current)
                .attr('width', width)
                .attr('height', height);

            svg.selectAll('*').remove();

            const g = svg.append('g').attr('class', 'graph');

            const zoom = d3.zoom()
                .scaleExtent([0.1, 4])
                .on('zoom', (event) => {
                    g.attr('transform', event.transform);
                });
// @ts-ignore
            svg.call(zoom);

            const simulation = d3.forceSimulation(nodes)
                // @ts-ignore
                .force('link', d3.forceLink().id(d => d.id).distance(150).strength(linkStrength))
                .force('charge', d3.forceManyBody().strength(-300))
                .force('center', d3.forceCenter(width / 2, height / 2));

            const link = g.append('g')
                .attr('class', 'edges')
                .selectAll('line')
                .data(links)
                .enter().append('line')
                .attr('stroke', linkStroke)
                .attr('stroke-opacity', linkStrokeOpacity)
                .attr('stroke-width', linkStrokeWidth)
                .attr('stroke-linecap', linkStrokeLinecap);

            const node = g.append('g')
                .attr('class', 'nodes')
                .selectAll('g')
                .data(nodes)
                .enter().append('g')
                // @ts-ignore
                .call(d3.drag()
                    .on('start', dragStarted)
                    .on('drag', dragged)
                );
                    // .on('end', dragEnded));

            node.append('circle')
                .attr('r', nodeRadius)
                // @ts-ignore
                .attr('fill', d => colors[d.group % colors.length])
                .attr('stroke', nodeStroke)
                .attr('stroke-width', nodeStrokeWidth)
                .attr('stroke-opacity', nodeStrokeOpacity);

            const labelGroup = g.append('g')
                .attr('class', 'labels');

            labelGroup.selectAll('text')
                .data(nodes)
                .enter().append('text')
                .attr('dy', -10)
                .attr('text-anchor', 'middle')
                // @ts-ignore
                .text(d => d.label);

            simulation
                .nodes(nodes)
                .on('tick', ticked);

            // @ts-ignore
            simulation.force('link')
                // @ts-ignore
                .links(links);

            function ticked() {
                link
                    // @ts-ignore
                    .attr('x1', d => d.source.x)
                    // @ts-ignore
                    .attr('y1', d => d.source.y)
                    // @ts-ignore
                    .attr('x2', d => d.target.x)
                    // @ts-ignore
                    .attr('y2', d => d.target.y);

                node
                    // @ts-ignore
                    .attr('transform', d => `translate(${d.x},${d.y})`);

                labelGroup.selectAll('text')
                    // @ts-ignore
                    .attr('transform', d => `translate(${d.x},${d.y})`);
            }

            function dragStarted(event: any, d: any) {
                if (!event.active) simulation.alphaTarget(0.3).restart();
                d.fx = d.x; // Fixar a posição no arrasto
                d.fy = d.y;
            }

            function dragged(event: any, d: any) {
                d.fx = event.x; // Atualiza a posição fixa
                d.fy = event.y;
            }

            // function dragEnded(event: { active: any; }, d: any) {
            //     if (!event.active) simulation.alphaTarget(0);
            //     // Não definir fx e fy como null para manter o nó fixo
            //     // d.fx = null;
            //     // d.fy = null;
            // }
        };

        drawGraph(data);
    }, [data, linkStroke, linkStrokeOpacity, linkStrokeWidth, linkStrokeLinecap, linkStrength, colors, nodeStroke, nodeStrokeWidth, nodeStrokeOpacity, nodeRadius]);

    return <svg ref={svgRef}></svg>;
}


export function Graph2({
                           data,
                           linkStroke = "#999",
                           linkStrokeOpacity = 0.6,
                           linkStrokeWidth = 1.5,
                           linkStrokeLinecap = "round",
                           linkStrength = 1,
                           colors = d3.schemeTableau10,
                           nodeStroke = "#fff",
                           nodeStrokeWidth = 1.5,
                           nodeStrokeOpacity = 1,
                           nodeRadius = 5,
                       }: GraphD3Props) {
    const svgRef = useRef<SVGSVGElement | null>(null);

    useEffect(() => {
        const drawGraph = (data: DataResponse) => {
            const links: LinkProps[] = data.edges.map(edge => ({
                source: edge.from,
                target: edge.to,
            }));
            const nodes: NodeProps[] = data.nodes.map(node => ({
                id: node.id,
                x: 0,
                y: 0,
                label: node.label,
                group: node.group || 0,
            }));

            const width = 900;
            const height = 800;

            const svg = d3.select(svgRef.current)
                .attr('width', width)
                .attr('height', height);

            svg.selectAll('*').remove();

            const g = svg.append('g').attr('class', 'graph');

            const zoom = d3.zoom()
                .scaleExtent([0.1, 4])
                .on('zoom', (event) => {
                    g.attr('transform', event.transform);
                });

            // @ts-ignore
            svg.call(zoom);

            const simulation = d3.forceSimulation(nodes)
                // @ts-ignore
                .force('link', d3.forceLink().id(d => d.id).distance(150).strength(linkStrength))
                .force('charge', d3.forceManyBody().strength(-300))
                .force('center', d3.forceCenter(width / 2, height / 2));

            const link = g.append('g')
                .attr('class', 'edges')
                .selectAll('line')
                .data(links)
                .enter().append('line')
                .attr('stroke', linkStroke)
                .attr('stroke-opacity', linkStrokeOpacity)
                .attr('stroke-width', linkStrokeWidth)
                .attr('stroke-linecap', linkStrokeLinecap);

            // @ts-ignore
            const node = g.append('g')
                .attr('class', 'nodes')
                .selectAll('g')
                .data(nodes)
                .enter().append('g')

                // @ts-ignore
                .call(d3.drag()
                    .on('start', dragStarted)
                    .on('drag', dragged)
                );
                    // .on('end', dragEnded));

            // @ts-ignore
            const nodeColors = new Map();

            const circles = node.append('circle')
                .attr('r', nodeRadius)
                // @ts-ignore
                .attr('fill', d => nodeColors.get(d.id) || colors[d.group % colors.length])
                .attr('stroke', nodeStroke)
                .attr('stroke-width', nodeStrokeWidth)
                .attr('stroke-opacity', nodeStrokeOpacity);

            const labelGroup = g.append('g')
                .attr('class', 'labels');

            labelGroup.selectAll('text')
                .data(nodes)
                .enter().append('text')
                .attr('dy', -10)
                .attr('text-anchor', 'middle')
                // @ts-ignore
                .text(d => d.label);

            simulation
                .nodes(nodes)
                .on('tick', ticked);

            // @ts-ignore
            simulation.force('link')
                // @ts-ignore
                .links(links);

            function ticked() {
                link
                    // @ts-ignore
                    .attr('x1', d => d.source.x)
                    // @ts-ignore
                    .attr('y1', d => d.source.y)
                    // @ts-ignore
                    .attr('x2', d => d.target.x)
                    // @ts-ignore
                    .attr('y2', d => d.target.y);

                // @ts-ignore
                node.attr('transform', d => `translate(${d.x},${d.y})`);

                labelGroup.selectAll('text')
                    // @ts-ignore
                    .attr('transform', d => `translate(${d.x},${d.y})`);
            }

            function dragStarted(event:any, d:any) {
                if (!event.active) simulation.alphaTarget(0.3).restart();
                d.fx = d.x; // Fixar a posição no arrasto
                d.fy = d.y;

                d.group = Math.floor(Math.random() * 10)
            }

            function dragged(event:any, d:any) {
                d.fx = event.x; // Atualiza a posição fixa
                d.fy = event.y;



                nodeColors.set(d.id, d3.color(colors[d.group % colors.length])!.darker(1).toString());

                // @ts-ignore
                circles.attr('fill', n => nodeColors.get(n.id) || colors[n.group % colors.length]);

                // Atualiza a cor das arestas conectadas
                link.attr('stroke', l => {
                    // @ts-ignore
                    return (l.source.id === d.id || l.target.id === d.id)
                        // @ts-ignore
                        ? d3.color(linkStroke).darker(1).toString()
                        : linkStroke;
                });
            }

            // function dragEnded(event: any, d: any) {
            //     if (!event.active) simulation.alphaTarget(0);
            //     // Não redefinir fx e fy para manter a posição fixa
            // }
        };

        drawGraph(data);
    }, [data, linkStroke, linkStrokeOpacity, linkStrokeWidth, linkStrokeLinecap, linkStrength, colors, nodeStroke, nodeStrokeWidth, nodeStrokeOpacity, nodeRadius]);

    return <svg ref={svgRef}></svg>;
}