import { Flex } from '@mantine/core';
import {
	Bar,
	BarChart,
	CartesianGrid,
	ResponsiveContainer,
	Tooltip,
	XAxis,
	YAxis,
} from 'recharts';

interface StackedBarChartProps {
	data: Array<{
		label: string;
		value: number;
		stackValue?: number;
	}>;
}

// This is a hack to suppress the warning about missing defaultProps in recharts library as of version 2.12
// @link https://github.com/recharts/recharts/issues/3615
const error = console.error;
console.error = (...args: any) => {
	if (/defaultProps/.test(args[0])) return;
	error(...args);
};

export const StackedBarChart = ({ data }: StackedBarChartProps) => {
	return (
		<Flex h={400} my="lg">
			<ResponsiveContainer>
				<BarChart data={data}>
					<CartesianGrid />
					<XAxis dataKey="label" />
					<YAxis />
					<Tooltip />
					<Bar dataKey="value" name="Visitors" stackId="a" fill="#7D44C7" />
					<Bar
						dataKey="stackValue"
						name="Page Views"
						stackId="a"
						fill="#9D5DEF"
					/>
				</BarChart>
			</ResponsiveContainer>
		</Flex>
	);
};
