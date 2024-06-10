import { Flex } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';
import { format, parseISO } from 'date-fns';
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
		date: string;
		value: number;
		stackValue?: number;
	}>;
}

const error = console.error;
/* biome-ignore lint/suspicious/noExplicitAny: This is a hack to suppress the warning about missing defaultProps in recharts
library as of version 2.12. @link https://github.com/recharts/recharts/issues/3615 */
console.error = (...args: any) => {
	if (/defaultProps/.test(args[0])) return;
	error(...args);
};

export const StackedBarChart = ({ data }: StackedBarChartProps) => {
	const [searchParams] = useSearchParams();

	const intlFormatter = new Intl.DateTimeFormat('en', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
	});

	// eslint-disable-next-line unicorn/consistent-function-scoping
	let dateFormatter = (date: Date) => format(date, 'MMM, yyyy');
	const period = searchParams.get('period');
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
				<BarChart data={data}>
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
