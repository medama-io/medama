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
