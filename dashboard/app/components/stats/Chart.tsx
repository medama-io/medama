import {
	Bar,
	BarChart,
	CartesianGrid,
	Legend,
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
		<BarChart
			width={1000}
			height={500}
			data={data}
			margin={{
				top: 20,
				right: 30,
				left: 20,
				bottom: 5,
			}}
		>
			<CartesianGrid strokeDasharray="3 3" />
			<XAxis dataKey="label" />
			<YAxis />
			<Tooltip />
			<Legend />
			<Bar dataKey="value" name="Visitors" stackId="a" fill="#8884d8" />
			<Bar dataKey="stackValue" name="Page Views" stackId="a" fill="#82ca9d" />
		</BarChart>
	);
};
