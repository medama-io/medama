import {
	type AreaChartProps,
	type BarChartProps,
	AreaChart as MantineAreaChart,
	BarChart as MantineBarChart,
} from '@mantine/charts';
import { useSearchParams } from '@remix-run/react';
import { format, parseISO } from 'date-fns';
import React, { useMemo } from 'react';

import { Group } from '@/components/layout/Flex';

import { formatCount, formatDuration, formatPercentage } from './formatter';

import classes from './Chart.module.css';

interface ChartData {
	date: string;
	value: number;
}

interface ChartProps {
	label: string;
	type: 'bar' | 'area';
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
type FixedPeriod = 'today' | 'yesterday' | 'quarter' | 'half' | 'year' | 'all';
type Period = HourPeriod | DayPeriod | FixedPeriod;

const PERIODS = {
	TODAY: 'today',
	YESTERDAY: 'yesterday',
	QUARTER: 'quarter',
	HALF: 'half',
	YEAR: 'year',
	ALL: 'all',
} as const;

const intlFormatterBasic = new Intl.DateTimeFormat('en', {
	dateStyle: 'long',
});

const intlFormatterDay = new Intl.DateTimeFormat('en', {
	dateStyle: 'full',
});

const intlFormatterMonth = new Intl.DateTimeFormat('en', {
	month: 'long',
	year: 'numeric',
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

		const item = payload[0];
		if (!item) return null;

		const dateTimeFormat = useMemo(() => {
			if (
				period === null ||
				period === undefined ||
				period === PERIODS.TODAY ||
				period === PERIODS.YESTERDAY ||
				period?.endsWith('h')
			) {
				return intlFormatterAll;
			}

			if (period?.endsWith('d')) {
				return intlFormatterDay;
			}

			if (
				period === PERIODS.HALF ||
				period === PERIODS.YEAR ||
				period === PERIODS.ALL
			) {
				return intlFormatterMonth;
			}

			return intlFormatterBasic;
		}, [period]);

		const dateLabel = useMemo(() => {
			const value = dateTimeFormat.format(parseISO(date));

			if (period === PERIODS.QUARTER) {
				return `Week of ${value}`;
			}

			if (
				period === PERIODS.HALF ||
				period === PERIODS.YEAR ||
				period === PERIODS.ALL
			) {
				return `Month of ${value}`;
			}

			return value;
		}, [dateTimeFormat, date, period]);

		return (
			<div className={classes.tooltip}>
				<h3 className={classes.date}>{dateLabel}</h3>
				<Group>
					<Group style={{ gap: 8 }}>
						<div
							className={classes.swatch}
							style={{ backgroundColor: item.color }}
						/>
						<span>{label}</span>
					</Group>
					<span>{valueFormatter(item.value)}</span>
				</Group>
			</div>
		);
	},
);

const AreaChart = (props: AreaChartProps) => {
	return (
		<MantineAreaChart
			areaProps={{ radius: 8, isAnimationActive: true, animationDuration: 500 }}
			curveType="linear"
			{...props}
		/>
	);
};

const BarChart = (props: BarChartProps) => {
	return (
		<MantineBarChart
			barChartProps={{ barCategoryGap: '15%' }}
			barProps={{ radius: 8, isAnimationActive: true, maxBarSize: 50 }}
			{...props}
		/>
	);
};

const Chart = ({ type, label, data }: ChartProps) => {
	const [searchParams] = useSearchParams();
	const period = searchParams.get('period') as Period | null;

	const dateFormatter = useMemo(() => {
		if (
			period === null ||
			period === undefined ||
			period === PERIODS.TODAY ||
			period === PERIODS.YESTERDAY ||
			period?.endsWith('h')
		) {
			return (date: Date) => format(date, 'HH:mm');
		}

		if (period?.endsWith('d')) {
			if (Number.parseInt(period) <= 7) {
				return (date: Date) => format(date, 'EEEEEE, MMM d');
			}

			return (date: Date) => format(date, 'MMM d');
		}

		return (date: Date) => format(date, 'MMM, yyyy');
	}, [period]);

	const valueFormatter = useMemo(() => {
		if (label === 'Time Spent') return formatDuration;
		if (label === 'Bounce Rate') return formatPercentage;
		return formatCount;
	}, [label]);

	const chartStyleProps: BarChartProps & AreaChartProps = {
		h: 400,
		my: 'xl',
		data,
		dataKey: 'date',
		series: [{ name: 'value', label, color: '#9D5DEF' }],
		tickLine: 'y',
		xAxisProps: {
			tickFormatter: (value: string) => dateFormatter(parseISO(value)),
			minTickGap: 20,
			interval: 'equidistantPreserveStart',
		},
		valueFormatter,
		tooltipProps: {
			content: ({ label: date, payload }) => (
				<ChartTooltip
					label={label}
					date={date}
					period={period}
					payload={payload as TooltipPayload[]}
					valueFormatter={valueFormatter}
				/>
			),
		},
		strokeDasharray: 0,
	};

	if (type === 'bar') return <BarChart {...chartStyleProps} />;
	if (type === 'area') return <AreaChart {...chartStyleProps} />;
	return null;
};

export { Chart };
