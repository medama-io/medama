import { useSearchParams } from '@remix-run/react';
import React, { useMemo } from 'react';

import { Tooltip, TooltipProvider } from '@/components/Tooltip';
import { Group } from '@/components/layout/Flex';
import { useChartType } from '@/hooks/use-chart-type';

import { formatCount, formatDuration, formatPercentage } from './formatter';
import type { StatHeaderData } from './types';

import classes from './HeaderDataBox.module.css';

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
	} else if (stat.chart === 'duration') {
		changeValue = formatDuration(Number(changeValue));
	} else {
		changeValue = formatCount(changeValue);
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

	const change = useMemo(() => {
		if (isPercentage) return stat.current - (stat.previous ?? 0);
		return calculateChange(stat.current, stat.previous);
	}, [stat, isPercentage]);

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
		<TooltipProvider delayDuration={0}>
			<Tooltip content={tooltipLabel}>
				<button
					type="submit"
					className={classes.databox}
					data-active={isActive}
					aria-label={`${stat.label}: ${formattedValue}. ${tooltipLabel}`}
					tabIndex={0}
					onClick={handleClick}
					onKeyDown={(event) => {
						if (event.key === 'Enter') handleClick();
					}}
				>
					<span className={classes.value}>{formattedValue}</span>
					<Group
						style={{ gap: 12, marginTop: 8, justifyContent: 'flex-start' }}
					>
						<p className={classes.label}>{stat.label}</p>
						<div
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
							{isPercentage ? formatPercentage(change) : `${change}%`}
						</div>
					</Group>
				</button>
			</Tooltip>
		</TooltipProvider>
	);
});

export { HeaderDataBox };
