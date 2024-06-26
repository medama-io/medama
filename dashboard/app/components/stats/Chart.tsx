import { BarChart as MantineBarChart } from '@mantine/charts';
import { ColorSwatch, Group, Paper, Text } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';
import { format, parseISO } from 'date-fns';
import React, { useMemo } from 'react';

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
	period: Period | null;
	payload: TooltipPayload[];
	valueFormatter: (value: number) => string;
}

type HourPeriod = `${number}h`;
type DayPeriod = `${number}d`;
type FixedPeriod = 'today' | 'yesterday' | 'quarter';
type Period = HourPeriod | DayPeriod | FixedPeriod;

const PERIODS = {
	TODAY: 'today',
	YESTERDAY: 'yesterday',
	QUARTER: 'quarter',
} as const;

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

const ChartTooltip = React.memo(
	({ label, date, period, payload, valueFormatter }: ChartTooltipProps) => {
		if (!payload || !label || !date) return null;

		const dateTimeFormat = useMemo(() => {
			if (
				period === PERIODS.TODAY ||
				period === PERIODS.YESTERDAY ||
				period?.endsWith('h')
			) {
				return intlFormatterAll;
			}

			if (period?.endsWith('d')) {
				return intlFormatterDay;
			}

			return intlFormatterBasic;
		}, [period]);

		return (
			<Paper px="md" py="md" withBorder shadow="md" radius="md">
				<Text fw={500} mb="xs">
					{dateTimeFormat.format(parseISO(date))}
				</Text>
				{payload.map((item) => (
					<Group key={item.name} gap="xs" justify="space-between">
						<Group gap="sm">
							<ColorSwatch color={item.color} size={12} withShadow={false} />
							<Text fz="sm">{label}</Text>
						</Group>
						<Text fz="sm">{valueFormatter(item.value)}</Text>
					</Group>
				))}
			</Paper>
		);
	},
);

const BarChart = ({ label, data }: BarChartProps) => {
	const [searchParams] = useSearchParams();
	const period = searchParams.get('period') as Period | null;

	const dateFormatter = useMemo(() => {
		switch (true) {
			case period === null:
			case period === undefined:
			case period === PERIODS.TODAY:
			case period === PERIODS.YESTERDAY:
			case period?.endsWith('h'): {
				return (date: Date) => format(date, 'HH:mm');
			}
			case period?.endsWith('d') && Number.parseInt(period) <= 7: {
				return (date: Date) => format(date, 'EEEEEE, MMM d');
			}
			case period?.endsWith('d') && Number.parseInt(period) > 7:
			case period === PERIODS.QUARTER: {
				return (date: Date) => format(date, 'MMM d');
			}
		}

		return (date: Date) => format(date, 'MMM, yyyy');
	}, [period]);

	const valueFormatter = useMemo(() => {
		if (label === 'Time Spent') return formatDuration;
		if (label === 'Bounce Rate') return formatPercentage;
		return formatCount;
	}, [label]);

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
