import { Flex } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';
import { format, parseISO } from 'date-fns';
import {
	Bar,
	BarChart as ReBarChart,
	CartesianGrid,
	ResponsiveContainer,
	Tooltip,
	XAxis,
	YAxis,
} from 'recharts';

interface ChartData {
	date: string;
	value: number;
}

interface BarChartProps {
	label: string;
	data: ChartData[];
}

const intlFormatter = new Intl.DateTimeFormat('en', {
	year: 'numeric',
	month: 'short',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
});

const BarChart = ({ label, data }: BarChartProps) => {
	const [searchParams] = useSearchParams();
	const period = searchParams.get('period');

	let dateFormatter = (date: Date) => format(date, 'MMM, yyyy');
	switch (true) {
		case period === null:
		case period === undefined:
		case period === 'today':
		case period === 'yesterday':
		case period?.endsWith('h'): {
			dateFormatter = (date: Date) => format(date, 'HH:mm');
			break;
		}
		case period?.endsWith('d') && Number.parseInt(String(period)) <= 7: {
			dateFormatter = (date: Date) => format(date, 'EEEEEE, MMM d');
			break;
		}
		case period?.endsWith('d') && Number.parseInt(String(period)) > 7:
		case period === 'quarter': {
			dateFormatter = (date: Date) => format(date, 'MMM d');
			break;
		}
	}

	return (
		<Flex h={400} my="lg">
			<ResponsiveContainer>
				<ReBarChart data={data}>
					<CartesianGrid />
					<XAxis
						dataKey="date"
						interval="equidistantPreserveStart"
						tickFormatter={(value) => dateFormatter(parseISO(value))}
						minTickGap={20}
					/>
					<YAxis />
					<Tooltip
						labelFormatter={(value) => intlFormatter.format(parseISO(value))}
					/>
					<Bar dataKey="value" name={label} stackId="a" fill="#9D5DEF" />
				</ReBarChart>
			</ResponsiveContainer>
		</Flex>
	);
};

export { BarChart };
