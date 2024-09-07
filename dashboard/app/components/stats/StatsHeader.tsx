import { CalendarIcon } from '@radix-ui/react-icons';
import * as ToggleGroup from '@radix-ui/react-toggle-group';
import { useParams, useSearchParams } from '@remix-run/react';
import type React from 'react';
import { Fragment } from 'react';
import { ScrollContainer } from 'react-indiana-drag-scroll';

import { DatePickerRange, datePickerClasses } from '@/components/DatePicker';
import { DropdownSelect } from '@/components/DropdownSelect';
import { Tooltip, TooltipProvider } from '@/components/Tooltip';
import { IconAreaChart } from '@/components/icons/area';
import { IconBarChart } from '@/components/icons/bar';
import { Group } from '@/components/layout/Flex';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { useChartType } from '@/hooks/use-chart-type';
import { useDisclosure } from '@/hooks/use-disclosure';
import { useMediaQuery } from '@/hooks/use-media-query';

import { HeaderDataBox } from './HeaderDataBox';
import type { ChartType, StatHeaderData } from './types';

import classes from './StatsHeader.module.css';

interface StatsHeaderProps {
	stats: StatHeaderData[];
	chart: string;
	websites: string[];
}

const DATE_GROUP_END_VALUES = ['yesterday', '30d', 'year', 'all'];

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

const SegmentedChartControl = () => {
	const { setChartType, getChartType } = useChartType();
	const chartType = getChartType();

	const handleChartChange = (value: ChartType) => {
		setChartType(value);
	};

	const chartTypes = CHART_TYPES.map((item) => (
		<Fragment key={item.value}>
			<Tooltip content={item.label}>
				<ToggleGroup.Item value={item.value} asChild>
					<button
						type="submit"
						className={classes.control}
						aria-label={item.label}
						onClick={() => handleChartChange(item.value)}
						onKeyDown={(e) => {
							if (e.key === 'Enter') {
								handleChartChange(item.value);
							}
						}}
						data-active={chartType === item.value}
					>
						<span className={classes.controlLabel}>{item.icon}</span>
					</button>
				</ToggleGroup.Item>
			</Tooltip>
		</Fragment>
	));

	return (
		<ToggleGroup.Root
			type="single"
			value={chartType}
			onValueChange={handleChartChange}
			asChild
		>
			<div className={classes.toggle}>
				<TooltipProvider delayDuration={500}>{chartTypes}</TooltipProvider>
			</div>
		</ToggleGroup.Root>
	);
};

const StatsHeader = ({ stats, chart, websites }: StatsHeaderProps) => {
	const { hostname } = useParams();
	const [searchParams] = useSearchParams();
	const [dateOpened, { toggle: toggleDate }] = useDisclosure(false);
	const hideWebsiteSelector = useMediaQuery('(36em < width < 62em)');

	// Convert websites array to object with same key-val for DropdownSelect
	const websitesRecord = Object.fromEntries(
		websites.map((website) => [website, website]),
	);

	const DATE_PRESETS: Record<string, React.ReactNode> = {
		today: 'Today',
		yesterday: 'Yesterday',
		'12h': 'Previous 12 hours',
		'24h': 'Previous 24 hours',
		'72h': 'Previous 72 hours',
		'7d': 'Previous 7 days',
		'14d': 'Previous 14 days',
		'30d': 'Previous 30 days',
		quarter: 'Previous quarter',
		half: 'Previous half year',
		year: 'Previous year',
		all: 'All time',
		custom: (
			<button
				type="button"
				onClick={toggleDate}
				className={datePickerClasses.button}
				data-state={
					searchParams.get('period') === 'custom' ? 'checked' : 'unchecked'
				}
			>
				Custom range
			</button>
		),
	};

	return (
		<InnerHeader>
			<div className={classes.title}>
				<h1>Dashboard</h1>
				<div className={classes.dropdowns}>
					{!hideWebsiteSelector && (
						<DropdownSelect
							records={websitesRecord}
							defaultValue={hostname ?? ''}
							defaultLabel="Unknown"
							ariaLabel="Select website"
							shouldUseNavigate
						/>
					)}
					<DropdownSelect
						records={DATE_PRESETS}
						defaultValue="today"
						defaultLabel="Custom range"
						ariaLabel="Select date range"
						icon={CalendarIcon}
						searchParamKey="period"
						separatorValues={DATE_GROUP_END_VALUES}
					/>
					<DatePickerRange open={dateOpened} setOpen={toggleDate} />
				</div>
			</div>
			<ScrollContainer className={classes.scrollcontainer}>
				<Group
					style={{
						alignItems: 'flex-end',
						marginTop: 8,
					}}
				>
					<div className={classes.scrollgroup}>
						{stats.map((stat) => (
							<HeaderDataBox
								key={stat.label}
								stat={stat}
								isActive={chart === stat.chart}
							/>
						))}
					</div>
					<SegmentedChartControl />
				</Group>
			</ScrollContainer>
		</InnerHeader>
	);
};

export { StatsHeader };
