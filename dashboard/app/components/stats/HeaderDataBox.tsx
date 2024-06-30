import { Box, Group, Tooltip, UnstyledButton } from '@mantine/core';
import React, { useMemo } from 'react';

import { useChartType } from '@/hooks/use-chart-type';

import { formatCount, formatDuration, formatPercentage } from './formatter';
import type { StatHeaderData } from './types';

import classes from './HeaderDataBox.module.css';
import { useSearchParams } from '@remix-run/react';

interface HeaderDataBoxProps {
	stat: StatHeaderData;
	isActive: boolean;
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
	isToday: boolean,
): string => {
	const period = isToday ? 'yesterday' : 'previous period';
	if (stat.previous === undefined || status === 'zero')
		return `No change since ${period}.`;

	let changeValue: string | number = Math.abs(stat.current - stat.previous);
	if (stat.chart === 'bounces') {
		changeValue = formatPercentage(changeValue);
	}

	if (stat.chart === 'duration') {
		changeValue = formatDuration(Number(changeValue));
	}

	return status === 'positive'
		? `Increased by ${changeValue} since ${period}.`
		: `Decreased by ${changeValue} since ${period}.`;
};

const HeaderDataBox = React.memo(({ stat, isActive }: HeaderDataBoxProps) => {
	const { setChartStat } = useChartType();
	const [searchParams] = useSearchParams();
	const period = searchParams.get('period') as 'today' | '24h' | string | null;
	const isToday = period === 'today' || period === '24h';

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
		() => formatTooltipLabel(stat, status, isToday),
		[stat, status, isToday],
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

export { HeaderDataBox };
