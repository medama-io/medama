import { Box, Flex, Group, Text, Tooltip } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';

import { type DataResponse } from '@/api/client';

import { DateComboBox } from './DateSelector';
import { formatCount, formatDuration, formatPercentage } from './formatter';
import classes from './StatsHeader.module.css';

interface HeaderDataBoxProps {
	label: string;
	value: number;
	previousValue?: number;
	isBounce?: boolean;
	isDuration?: boolean;
	isActive?: boolean;
	hideBadge?: boolean;
}

const HeaderDataBox = ({
	label,
	value,
	previousValue,
	isBounce,
	isDuration,
	isActive,
	hideBadge,
}: HeaderDataBoxProps) => {
	// Calculate percentage change if previous value is available
	let change = 0;
	if (previousValue) {
		// If isBounce, it is already a percentage so we just need to calculate
		// the difference between the current and previous bounce rates
		change = isBounce
			? Math.round((value - previousValue) * 100)
			: Math.round(((value - previousValue) / previousValue) * 100);
	}

	let formattedValue: string;
	if (isDuration) {
		// Format the duration into a human readable format
		formattedValue = formatDuration(value);
	} else {
		// Rely on Intl.NumberFormat to format the values according to the user's locale
		formattedValue = isBounce ? formatPercentage(value) : formatCount(value);
	}

	// Determine if the change is positive or negative
	let status: 'positive' | 'negative' | 'zero' = 'zero';
	if (change > 0) {
		status = 'positive';
	} else if (change < 0) {
		status = 'negative';
	}

	let badgeColor: string;
	if (isBounce) {
		badgeColor = status === 'positive' ? '#FFD5B7' : '#DFFFB7';
	} else {
		badgeColor = status === 'positive' ? '#DFFFB7' : '#FFD5B7';
	}

	// Generate a tooltip label depending on if the change is positive or negative
	let tooltipLabel = 'No change since yesterday.';
	if (previousValue !== undefined && status !== 'zero') {
		let changeValue: number | string;
		if (isBounce) {
			changeValue = `${Math.round(Math.abs(value - previousValue) * 100)}%`;
		} else if (isDuration) {
			changeValue = formatDuration(Math.abs(value - previousValue));
		} else {
			changeValue = Math.abs(value - previousValue);
		}

		if (status === 'positive') {
			tooltipLabel = `Increased by ${changeValue} since yesterday.`;
		} else if (status === 'negative') {
			tooltipLabel = `Decreased by ${changeValue} since yesterday.`;
		}
	}

	return (
		<Tooltip label={tooltipLabel} withArrow>
			<Box
				className={classes['data-box']}
				bg={isActive ? '#39414E' : undefined}
			>
				<Text fw={600} fz={28} pb={6}>
					{formattedValue}
				</Text>
				<Group gap="xs">
					<Text fz={14} span>
						{label}
					</Text>
					{!hideBadge && (
						<Box className={classes.badge} bg={badgeColor}>
							{status === 'positive' ? '+' : undefined}
							{change}%
						</Box>
					)}
				</Group>
			</Box>
		</Tooltip>
	);
};

type StatsHeaderProps = NonNullable<DataResponse<'StatsSummary'>['data']>;

export const StatsHeader = ({ current, previous }: StatsHeaderProps) => {
	// Calculate current bounce rate by dividing the number of bounces to the total number of unique visitors
	const bounceRate = current.bounces / current.visitors || 0; // Avoid NaN
	const previousBounceRate = previous
		? previous.bounces / previous.visitors || 0
		: 0;

	const [searchParams] = useSearchParams();
	const isAllTime = searchParams.get('period') === 'all';

	return (
		<div className={classes.header}>
			<div className={classes.inner}>
				<Flex justify="space-between">
					<Text fw={500} fz={32} pb="xl">
						Dashboard
					</Text>
					<DateComboBox />
				</Flex>
				<Group>
					<Group>
						<HeaderDataBox
							label="Visitors"
							value={current.visitors}
							previousValue={previous?.visitors}
							hideBadge={isAllTime}
							isActive
						/>
						<HeaderDataBox
							label="Page Views"
							value={current.pageviews}
							previousValue={previous?.pageviews}
							hideBadge={isAllTime}
						/>
						<HeaderDataBox
							label="Time Spent"
							value={current.duration}
							previousValue={previous?.duration}
							hideBadge={isAllTime}
							isDuration
						/>
						<HeaderDataBox
							label="Bounce Rate"
							value={bounceRate}
							previousValue={previousBounceRate}
							hideBadge={isAllTime}
							isBounce
						/>
					</Group>
					<div>Switch Chart</div>
				</Group>
			</div>
		</div>
	);
};
