import { Box, Flex, Group, Tooltip, UnstyledButton } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';
import React, { useMemo } from 'react';

import { InnerHeader } from '@/components/layout/InnerHeader';

import { DateComboBox } from './DateSelector';
import { formatCount, formatDuration, formatPercentage } from './formatter';
import type { StatHeaderData } from './types';

import classes from './StatsHeader.module.css';

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
): string => {
	if (stat.previous === undefined || status === 'zero')
		return 'No change since yesterday.';

	const isPercentage = stat.chart === 'bounces';
	const isDuration = stat.chart === 'duration';

	// Rely on Intl.NumberFormat to format the values according to the user's locale
	const changeValue = isPercentage
		? `${Math.abs(stat.current - stat.previous).toFixed(2)}%`
		: isDuration
			? formatDuration(Math.abs(stat.current - stat.previous))
			: Math.abs(stat.current - stat.previous);

	return status === 'positive'
		? `Increased by ${changeValue} since yesterday.`
		: `Decreased by ${changeValue} since yesterday.`;
};

const HeaderDataBox = React.memo(({ stat, isActive }: HeaderDataBoxProps) => {
	const [searchParams, setSearchParams] = useSearchParams();

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
		searchParams.set('chart', stat.chart);
		setSearchParams(searchParams);
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

interface StatsHeaderProps {
	stats: StatHeaderData[];
	chart: string;
}

const StatsHeader = ({ stats, chart }: StatsHeaderProps) => {
	return (
		<InnerHeader>
			<Flex justify="space-between" align="center" py={8}>
				<h1>Dashboard</h1>
				<DateComboBox />
			</Flex>
			<Group mt="xs">
				{stats.map((stat) => (
					<HeaderDataBox
						key={stat.label}
						stat={stat}
						isActive={chart === stat.chart}
					/>
				))}
			</Group>
		</InnerHeader>
	);
};

export { StatsHeader };
