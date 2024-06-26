import { BarChart as MantineBarChart } from '@mantine/charts';
import { useSearchParams } from '@remix-run/react';
import { format, parseISO } from 'date-fns';
import { ColorSwatch, Group, Paper, Text } from '@mantine/core';

import { formatCount, formatDuration, formatPercentage } from './formatter';

interface ChartData {
	date: string;
	value: number;
}

interface BarChartProps {
	label: string;
	data: ChartData[];
}

interface TooltipPayload {
	name: string;
	value: number;
	color: string;
}

interface ChartTooltipProps {
	label: string;
	date: string;
	period: string | null;
	payload: TooltipPayload[];
	valueFormatter: (value: number) => string;
}

const intlFormatterBasic = new Intl.DateTimeFormat('en', {
	dateStyle: 'long',
});

const intlFormatterDay = new Intl.DateTimeFormat('en', {
	dateStyle: 'full',
});

const intlFormatterAll = new Intl.DateTimeFormat('en', {
	year: 'numeric',
	month: 'short',
	day: 'numeric',
	hour: 'numeric',
});

const ChartTooltip = ({
	label,
	date,
	period,
	payload,
	valueFormatter,
}: ChartTooltipProps) => {
	if (!payload || !label || !date) return null;

	let dateFormatter = intlFormatterBasic;
	if (period === 'today' || period === 'yesterday' || period?.endsWith('h')) {
		dateFormatter = intlFormatterAll;
	}
	if (period?.endsWith('d')) {
		dateFormatter = intlFormatterDay;
	}

	return (
		<Paper px="md" py="md" withBorder shadow="md" radius="md">
			<Text fw={500} mb="xs">
				{dateFormatter.format(parseISO(date))}
			</Text>
			{payload.map((item) => (
				<Group key={item.name} gap="xs" justify="space-between">
					<Group gap="sm">
						<ColorSwatch color={item.color} size={12} withShadow={false} />
						<Text key={item.name} fz="sm">
							{label}
						</Text>
					</Group>
					<Text key={item.name} fz="sm">
						{valueFormatter(item.value)}
					</Text>
				</Group>
			))}
		</Paper>
	);
};

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

	const isDuration = label === 'Time Spent';
	const isPercentage = label === 'Bounce Rate';
	const valueFormatter = isDuration
		? formatDuration
		: isPercentage
			? formatPercentage
			: formatCount;

	return (
		<MantineBarChart
			h={400}
			my="xl"
			data={data}
			dataKey="date"
			series={[{ name: 'value', label: 'Visitors', color: '#9D5DEF' }]}
			barProps={{ radius: 8, isAnimationActive: true }}
			tickLine="y"
			xAxisProps={{
				tickFormatter: (value) => dateFormatter(parseISO(value)),
				minTickGap: 20,
				interval: 'equidistantPreserveStart',
			}}
			valueFormatter={valueFormatter}
			tooltipProps={{
				content: ({ label: date, payload }) => (
					<ChartTooltip
						label={label}
						date={date}
						period={period}
						payload={payload as TooltipPayload[]}
						valueFormatter={valueFormatter}
					/>
				),
			}}
			strokeDasharray={0}
		/>
	);
};

export { BarChart };
