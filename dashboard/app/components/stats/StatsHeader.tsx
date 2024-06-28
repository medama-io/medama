import {
	Box,
	Flex,
	FloatingIndicator,
	Group,
	Tooltip,
	UnstyledButton,
} from '@mantine/core';
import React, { useMemo, useState } from 'react';

import { IconAreaChart } from '@/components/icons/area';
import { IconBarChart } from '@/components/icons/bar';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { useChartType } from '@/hooks/use-chart-type';

import { DateComboBox } from './DateSelector';
import { formatCount, formatDuration, formatPercentage } from './formatter';
import type { ChartType, StatHeaderData } from './types';

import classes from './StatsHeader.module.css';

interface HeaderDataBoxProps {
	stat: StatHeaderData;
	isActive: boolean;
}

interface StatsHeaderProps {
	stats: StatHeaderData[];
	chart: string;
}

// Calculate percentage change if previous value is available.
const calculateChange = (current: number, previous?: number): number => {
	if (previous) {
		return Math.round(((current - previous) / previous) * 100);
	}
	return 0;
};

const getStatus = (change: number): 'positive' | 'negative' | 'zero' => {
	if (change > 0) return 'positive';
	if (change < 0) return 'negative';
	return 'zero';
};

const formatTooltipLabel = (
	stat: StatHeaderData,
	status: 'positive' | 'negative' | 'zero',
): string => {
	if (stat.previous === undefined || status === 'zero')
		return 'No change since yesterday.';

	let changeValue: string | number = Math.abs(stat.current - stat.previous);
	if (stat.chart === 'bounces') {
		changeValue = formatPercentage(changeValue);
	}

	if (stat.chart === 'duration') {
		changeValue = formatDuration(Number(changeValue));
	}

	return status === 'positive'
		? `Increased by ${changeValue} since yesterday.`
		: `Decreased by ${changeValue} since yesterday.`;
};

const HeaderDataBox = React.memo(({ stat, isActive }: HeaderDataBoxProps) => {
	const { setChartStat } = useChartType();

	const isPercentage = stat.chart === 'bounces';
	const isDuration = stat.chart === 'duration';

	const change = useMemo(
		() => calculateChange(stat.current, stat.previous),
		[stat],
	);

	const status = useMemo(() => getStatus(change), [change]);
	const formattedValue = useMemo(
		() =>
			isDuration
				? formatDuration(stat.current)
				: isPercentage
					? formatPercentage(stat.current)
					: formatCount(stat.current),
		[stat, isDuration, isPercentage],
	);

	const tooltipLabel = useMemo(
		() => formatTooltipLabel(stat, status),
		[stat, status],
	);

	const handleClick = () => {
		setChartStat(stat.chart);
	};

	return (
		<Tooltip label={tooltipLabel} withArrow>
			<UnstyledButton
				className={classes.databox}
				data-active={isActive}
				aria-label={`${stat.label}: ${formattedValue}. ${tooltipLabel}`}
				role="region"
				tabIndex={0}
				onClick={handleClick}
			>
				<span className={classes.value}>{formattedValue}</span>
				<Group gap="sm" mt={8}>
					<p className={classes.label}>{stat.label}</p>
					<Box
						className={classes.badge}
						data-status={
							isPercentage
								? status === 'positive'
									? 'negative'
									: 'positive'
								: status
						}
						aria-label={`Change: ${change}%`}
						role="status"
					>
						{status === 'positive' ? '+' : ''}
						{change}%
					</Box>
				</Group>
			</UnstyledButton>
		</Tooltip>
	);
});

const CHART_TYPES = [
	{
		label: 'Toggle area chart',
		value: 'area',
		icon: <IconAreaChart />,
	},
	{
		label: 'Toggle bar chart',
		value: 'bar',
		icon: <IconBarChart />,
	},
] as const;

const StatsHeader = ({ stats, chart }: StatsHeaderProps) => {
	const [rootRef, setRootRef] = useState<HTMLDivElement | null>(null);
	const [controlsRefs, setControlsRefs] = useState<
		Record<string, HTMLButtonElement | null>
	>({});

	const { setChartType, getChartType } = useChartType();
	const chartType = getChartType();

	const handleChartChange = (value: ChartType) => {
		setChartType(value);
	};

	const setControlRef = (type: ChartType) => (node: HTMLButtonElement) => {
		controlsRefs[type] = node;
		setControlsRefs(controlsRefs);
	};

	const chartTypes = CHART_TYPES.map((item) => (
		<Tooltip key={item.value} label={item.label} withArrow>
			<UnstyledButton
				key={item.value}
				className={classes.control}
				ref={setControlRef(item.value)}
				aria-label={item.label}
				onClick={() => handleChartChange(item.value)}
				data-active={chartType === item.value}
			>
				<span className={classes.controlLabel}>{item.icon}</span>
			</UnstyledButton>
		</Tooltip>
	));

	return (
		<InnerHeader>
			<Flex justify="space-between" align="center" py={8}>
				<h1>Dashboard</h1>
				<DateComboBox />
			</Flex>
			<Group justify="space-between" align="flex-end">
				<Group mt="xs">
					{stats.map((stat) => (
						<HeaderDataBox
							key={stat.label}
							stat={stat}
							isActive={chart === stat.chart}
						/>
					))}
				</Group>
				<div className={classes.toggle} ref={setRootRef}>
					<Tooltip.Group openDelay={1000}>{chartTypes}</Tooltip.Group>
					<FloatingIndicator
						component="span"
						className={classes.indicator}
						target={controlsRefs[chartType]}
						parent={rootRef}
					/>
				</div>
			</Group>
		</InnerHeader>
	);
};

export { StatsHeader };
