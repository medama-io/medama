import { Box, Flex, Group, Tooltip, UnstyledButton } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';

import type { DataResponse } from '@/api/client';

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

// Calculate percentage change if previous value is available.
const calculateChange = (
	value: number,
	previousValue?: number,
	isBounce?: boolean,
): number => {
	// If isBounce, it is already a percentage so we just need to calculate
	// the difference between the current and previous bounce rates.
	if (previousValue) {
		return isBounce
			? Math.round((value - previousValue) * 100)
			: Math.round(((value - previousValue) / previousValue) * 100);
	}
	return 0;
};

const getStatus = (change: number): 'positive' | 'negative' | 'zero' => {
	if (change > 0) return 'positive';
	if (change < 0) return 'negative';
	return 'zero';
};

const formatTooltipLabel = (
	value: number,
	previousValue: number | undefined,
	status: 'positive' | 'negative' | 'zero',
	isBounce: boolean | undefined,
	isDuration: boolean | undefined,
): string => {
	if (previousValue === undefined || status === 'zero') {
		return 'No change since yesterday.';
	}

	// Rely on Intl.NumberFormat to format the values according to the user's locale
	const changeValue = isBounce
		? `${Math.round(Math.abs(value - previousValue) * 100)}%`
		: isDuration
			? formatDuration(Math.abs(value - previousValue))
			: Math.abs(value - previousValue);

	return status === 'positive'
		? `Increased by ${changeValue} since yesterday.`
		: `Decreased by ${changeValue} since yesterday.`;
};

const HeaderDataBox = ({
	label,
	value,
	previousValue,
	isBounce,
	isDuration,
	isActive,
	hideBadge,
}: HeaderDataBoxProps) => {
	const change = calculateChange(value, previousValue, isBounce);
	const status = getStatus(change);
	const formattedValue = isDuration
		? formatDuration(value)
		: isBounce
			? formatPercentage(value)
			: formatCount(value);

	const tooltipLabel = formatTooltipLabel(
		value,
		previousValue,
		status,
		isBounce,
		isDuration,
	);

	return (
		<Tooltip label={tooltipLabel} withArrow>
			<UnstyledButton
				className={classes.databox}
				data-active={isActive}
				aria-label={`${label} is ${formattedValue}. ${tooltipLabel}`}
				role="region"
				tabIndex={0}
			>
				<span className={classes.value}>{formattedValue}</span>
				<Group gap="sm" mt={8}>
					<p className={classes.label}>{label}</p>
					{!hideBadge && (
						<Box
							className={classes.badge}
							data-status={
								isBounce
									? status === 'positive'
										? 'negative'
										: 'positive'
									: status
							}
							aria-label={`Difference is ${change}%`}
							role="status"
						>
							{status === 'positive' ? '+' : undefined}
							{change}%
						</Box>
					)}
				</Group>
			</UnstyledButton>
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
				<Flex justify="space-between" align="center" py={8}>
					<p className={classes['header-title']} role="heading" aria-level={1}>
						Dashboard
					</p>
					<DateComboBox />
				</Flex>
				<Group mt="xs">
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
			</div>
		</div>
	);
};
