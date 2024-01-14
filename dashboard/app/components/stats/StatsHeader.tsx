import { Box, Flex, Group, Text, Tooltip } from '@mantine/core';

import { type DataResponse } from '@/api/client';

import classes from './StatsHeader.module.css';

// Convert a duration in milliseconds to a human readable format
// such as 1m30s or 30s
const formatDuration = (durationMs: number) => {
	const totalSeconds = Math.floor(durationMs / 1000);
	const minutes = Math.floor(totalSeconds / 60);
	const seconds = totalSeconds % 60;

	return minutes === 0 ? `${seconds}s` : `${minutes}m${seconds}s`;
};

interface HeaderDataBoxProps {
	label: string;
	value: number;
	previousValue?: number;
	isBounce?: boolean;
	isDuration?: boolean;
	isLive?: boolean;
	isActive?: boolean;
}

const HeaderDataBox = ({
	label,
	value,
	previousValue,
	isBounce,
	isDuration,
	isLive,
	isActive,
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

	// Format values into a more readable format
	// navigator.languages has a readonly property so we need to copy it into a new array
	const languages =
		typeof document === 'undefined' ? ['en-US'] : [...navigator.languages];

	let formattedValue: string;
	if (isDuration) {
		// Format the duration into a human readable format
		formattedValue = formatDuration(value);
	} else {
		// Rely on Intl.NumberFormat to format the values according to the user's locale
		const formatter = isBounce
			? Intl.NumberFormat(languages, {
					style: 'percent',
					maximumFractionDigits: 1,
			  })
			: Intl.NumberFormat(languages, {
					notation: 'compact',
					maximumFractionDigits: 2,
			  });

		formattedValue = formatter.format(value);
	}

	// Determine if the change is positive or negative
	let status: 'positive' | 'negative' | 'zero' = 'zero';
	if (change > 0) {
		status = 'positive';
	} else if (change < 0) {
		status = 'negative';
	}

	// Generate a tooltip label depending on if the change is positive or negative
	let tooltipLabel = 'No change since yesterday.';
	if (previousValue !== undefined) {
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
		<Tooltip label={tooltipLabel} disabled={isLive}>
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
					{change !== undefined && !isLive && (
						<Box
							className={classes.badge}
							bg={
								status === 'positive' || status === 'zero'
									? '#DFFFB7'
									: '#FFD5B7'
							}
						>
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
	const bounceRate = current.bounces / current.uniques;
	const previousBounceRate = previous ? previous.bounces / previous.uniques : 0;

	return (
		<div className={classes.header}>
			<div className={classes.inner}>
				<Flex justify="space-between">
					<Text fw={500} fz={32} pb="xl">
						Dashboard
					</Text>
					<Flex>Selector</Flex>
				</Flex>
				<Group>
					<Group>
						<HeaderDataBox
							label="Unique Visitors"
							value={current.uniques}
							previousValue={previous?.uniques}
							isActive
						/>
						<HeaderDataBox
							label="Page Views"
							value={current.pageviews}
							previousValue={previous?.pageviews}
						/>
						<HeaderDataBox
							label="Avg. Duration"
							value={current.duration}
							previousValue={previous?.duration}
							isDuration
						/>
						<HeaderDataBox
							label="Bounce Rate"
							value={bounceRate}
							previousValue={previousBounceRate}
							isBounce
						/>
						<HeaderDataBox label="Active" value={current.active} isLive />
					</Group>
					<div>Switch Chart</div>
				</Group>
			</div>
		</div>
	);
};
