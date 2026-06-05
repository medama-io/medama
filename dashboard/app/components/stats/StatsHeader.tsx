import { SegmentedControl, VisuallyHidden } from '@mantine/core';
import { useDisclosure, useMediaQuery } from '@mantine/hooks';
import { Calendar } from 'lucide-react';
import type React from 'react';
import { ScrollContainer } from 'react-indiana-drag-scroll';
import { useParams } from 'react-router';

import { DatePickerRange } from '@/components/DatePicker';
import { DropdownSelect } from '@/components/DropdownSelect';
import { IconAreaChart } from '@/components/icons/area';
import { IconBarChart } from '@/components/icons/bar';
import { Group } from '@/components/layout/Flex';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { Tooltip, TooltipProvider } from '@/components/Tooltip';
import { useChartType } from '@/hooks/use-chart-type';

import { HeaderDataBox } from './HeaderDataBox';
import classes from './StatsHeader.module.css';
import type { ChartType, StatHeaderData } from './types';

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

	const chartTypes = CHART_TYPES.map((item) => ({
		value: item.value,
		label: (
			<Tooltip content={item.label}>
				<span className={classes.controlIcon} data-chart-toggle={item.value}>
					{item.icon}
					<VisuallyHidden>{item.label}</VisuallyHidden>
				</span>
			</Tooltip>
		),
	}));

	return (
		<TooltipProvider delayDuration={500}>
			<SegmentedControl
				value={chartType}
				onChange={(value) => handleChartChange(value as ChartType)}
				data={chartTypes}
				classNames={{
					root: classes.toggle,
					indicator: classes.indicator,
					control: classes.control,
					label: classes.controlLabel,
				}}
			/>
		</TooltipProvider>
	);
};

const StatsHeader = ({ stats, chart, websites }: StatsHeaderProps) => {
	const { hostname } = useParams();
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
		custom: 'Custom range',
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
						icon={Calendar}
						customActions={{ custom: toggleDate }}
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
