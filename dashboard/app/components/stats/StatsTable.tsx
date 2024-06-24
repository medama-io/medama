import { ActionIcon, Box, Group, Tabs, Text } from '@mantine/core';
import { Link, useNavigate, useSearchParams } from '@remix-run/react';
import {
	DataTable,
	type DataTableColumn,
	type DataTableSortStatus,
} from 'mantine-datatable';
import { useEffect, useState } from 'react';

import { IconChevronLeft } from '@/components/icons/chevronleft';
import { IconChevronRight } from '@/components/icons/chevronright';

import { formatCount, formatDuration, formatPercentage } from './formatter';
import classes from './StatsTable.module.css';

interface DataRow {
	// Common
	visitors?: number;
	visitors_percentage?: number;
	// Mixed
	path?: string;
	bounces?: number;
	bounce_rate?: number;
	duration?: number;
	// Pages
	pageviews?: number;
	pageviews_percentage?: number;
	// Duration
	duration_upper_quartile?: number;
	duration_lower_quartile?: number;
	duration_percentage?: number;
	// Referrers
	referrer?: string;
	// Sources
	source?: string;
	// Mediums
	medium?: string;
	// Campaigns
	campaign?: string;
	// Browsers
	browser?: string;
	// Operating Systems
	os?: string;
	// Devices
	device?: string;
	// Countries
	country?: string;
	// Languages
	language?: string;
}

const labelMap = {
	pages: 'Pages',
	time: 'Time Spent',
	referrers: 'Referrers',
	sources: 'Sources',
	mediums: 'Mediums',
	campaigns: 'Campaigns',
	browsers: 'Browsers',
	os: 'Operating Systems',
	devices: 'Devices',
	countries: 'Countries',
	languages: 'Languages',
};

interface QueryTableProps {
	query: keyof typeof labelMap | string;
	data: DataRow[];
}

// Preset columns
const path = { accessor: 'path', title: 'Path', width: '100%' };
const visitors: DataTableColumn = {
	accessor: 'visitors',
	title: 'Visitors',
	sortable: true,
};
const visitorsPercentage: DataTableColumn = {
	accessor: 'visitors_percentage',
	title: 'Visitors %',
	render: (record: DataRow) =>
		formatPercentage(record.visitors_percentage ?? 0),
};
const pageviews: DataTableColumn = {
	accessor: 'pageviews',
	title: 'Views',
	sortable: true,
	render: (record: DataRow) => formatCount(record.pageviews ?? 0),
};
const pageviewsPercentage: DataTableColumn = {
	accessor: 'pageviews_percentage',
	title: 'Views %',
	render: (record: DataRow) =>
		formatPercentage(record.pageviews_percentage ?? 0),
};
const bounceRate: DataTableColumn = {
	accessor: 'bounce_rate',
	title: 'Bounce %',
	render: (record: DataRow) =>
		formatPercentage((record.bounces ?? 0) / (record.visitors ?? 0)),
};
const duration: DataTableColumn = {
	accessor: 'duration',
	title: 'Duration',
	sortable: true,
	render: (record: DataRow) => formatDuration(record.duration ?? 0),
};

const PAGE_SIZES = [10, 25, 50, 100];

// biome-ignore lint/suspicious/noExplicitAny: Generic function
const sortBy = (key: any) => {
	// biome-ignore lint/suspicious/noExplicitAny: Generic function
	return (a: any, b: any) => (a[key] > b[key] ? 1 : b[key] > a[key] ? -1 : 0);
};

const QueryTable = ({ query, data }: QueryTableProps) => {
	// Pagination
	const [pageSize, setPageSize] = useState(10);
	const [page, setPage] = useState(1);
	const [records, setRecords] = useState(data.slice(0, pageSize));

	const handlePageChange = (newPage: number) => {
		// Prevent negative pages
		if (newPage < 1 || newPage > Math.ceil(data.length / pageSize)) {
			return;
		}

		setPage(newPage);
	};

	const handlePageSizeChange = (newSize: number) => {
		setPageSize(newSize);
		setPage(1);
	};

	// Sorting
	const [sortStatus, setSortStatus] = useState<DataTableSortStatus<DataRow>>({
		columnAccessor: 'visitors',
		direction: 'desc',
	});

	useEffect(() => {
		// Calculate the range of records to display
		const from = (page - 1) * pageSize;
		const to = from + pageSize;

		// Sort and slice the data to page size
		const temp = [...data].sort(sortBy(sortStatus.columnAccessor));
		setRecords(
			sortStatus.direction === 'desc'
				? temp.reverse().slice(from, to)
				: temp.slice(from, to),
		);
	}, [sortStatus, data, page, pageSize]);

	// Define columns based on query
	const columns: DataTableColumn[] = [];
	switch (query) {
		case 'pages': {
			columns.push(
				path,
				visitors,
				visitorsPercentage,
				pageviews,
				pageviewsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'time': {
			columns.push(
				path,
				visitors,
				duration,
				{
					accessor: 'duration_lower_quartile',
					title: 'Q1 (25%)',
					sortable: true,
					render: (record: DataRow) =>
						formatDuration(record.duration_lower_quartile ?? 0),
				},
				{
					accessor: 'duration_upper_quartile',
					title: 'Q3 (75%)',
					sortable: true,
					render: (record: DataRow) =>
						formatDuration(record.duration_upper_quartile ?? 0),
				},
				{
					accessor: 'duration_percentage',
					title: 'Duration %',
					sortable: true,
					render: (record: DataRow) =>
						formatPercentage((record.duration_percentage ?? 0) / 100),
				},
				bounceRate,
			);
			break;
		}
		case 'referrers': {
			columns.push(
				{
					accessor: 'referrer',
					title: 'Referrer',
					width: '100%',
					render: (record: DataRow) =>
						record.referrer === '' ? 'Direct/None' : record.referrer,
				},
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'sources': {
			columns.push(
				{
					accessor: 'source',
					title: 'Source',
					width: '100%',
					render: (record: DataRow) =>
						record.source === '' ? 'Direct/None' : record.source,
				},
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'mediums': {
			columns.push(
				{
					accessor: 'medium',
					title: 'Medium',
					width: '100%',
					render: (record: DataRow) =>
						record.medium === '' ? 'Direct/None' : record.medium,
				},
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'campaigns': {
			columns.push(
				{
					accessor: 'campaign',
					title: 'Campaign',
					width: '100%',
					render: (record: DataRow) =>
						record.campaign === '' ? 'Direct/None' : record.campaign,
				},
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'browsers': {
			columns.push(
				{ accessor: 'browser', title: 'Browser', width: '100%' },
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'os': {
			columns.push(
				{ accessor: 'os', title: 'Operating System', width: '100%' },
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'devices': {
			columns.push(
				{ accessor: 'device', title: 'Device', width: '100%' },
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'countries': {
			columns.push(
				{ accessor: 'country', title: 'Country', width: '100%' },
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		case 'languages': {
			columns.push(
				{ accessor: 'language', title: 'Language', width: '100%' },
				visitors,
				visitorsPercentage,
				bounceRate,
				duration,
			);
			break;
		}
		default: {
			return <div>Invalid query</div>;
		}
	}

	return (
		<div className={classes['table-wrapper']}>
			<div className={classes['table-header']}>
				<Text fz={14} fw={600} py={3}>
					{labelMap[query] ?? 'N/A'}
				</Text>
			</div>
			<DataTable
				minHeight={300}
				highlightOnHover
				withRowBorders={false}
				// Have to type assert here as technically we have Record<string | undefined, unknown>[]
				// but we don't know the exact shape of the data
				records={records as Array<Record<string, unknown>>}
				// biome-ignore lint/suspicious/noExplicitAny: Does not handle generic types very well
				columns={columns as any}
				sortStatus={sortStatus}
				onSortStatusChange={setSortStatus}
			/>
			<Group justify="space-between" px="lg" py="sm">
				<Group>
					<span className={classes.viewspan}>View</span>
					{PAGE_SIZES.map((size) => (
						<ActionIcon
							key={size}
							variant="transparent"
							className={classes['page-size']}
							onClick={() => {
								handlePageSizeChange(size);
							}}
							disabled={size === pageSize || data.length <= size}
							data-active={size === pageSize}
						>
							{size}
						</ActionIcon>
					))}
				</Group>
				<Group>
					<ActionIcon
						variant="transparent"
						className={classes['page-arrow']}
						onClick={() => {
							handlePageChange(page - 1);
						}}
						disabled={page <= 1}
					>
						<IconChevronLeft />
					</ActionIcon>
					<span>
						Page {page} of {Math.ceil(data.length / pageSize)}
					</span>
					<ActionIcon
						variant="transparent"
						className={classes['page-arrow']}
						onClick={() => {
							handlePageChange(page + 1);
						}}
						disabled={page >= Math.ceil(data.length / pageSize)}
					>
						<IconChevronRight />
					</ActionIcon>
				</Group>
			</Group>
		</div>
	);
};

interface StatsTableProps {
	query: string;
	data: DataRow[];
}

export const StatsTable = ({ query, data }: StatsTableProps) => {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();

	return (
		<Tabs
			variant="unstyled"
			value={query}
			classNames={{
				root: classes['tab-root'],
				tab: classes.tab,
				list: classes['tab-list'],
			}}
			orientation="vertical"
			onChange={(value) => {
				navigate(
					{
						pathname: `../${value}`,
						// Preserve search params when switching tabs
						search: `?${searchParams.toString()}`,
					},
					{ preventScrollReset: true },
				);
			}}
			keepMounted={false}
		>
			<Tabs.List>
				<div className={classes['list-wrapper']}>
					<Box
						component={Link}
						to={{
							pathname: '../',
							search: `?${searchParams.toString()}`,
						}}
						className={classes.back}
					>
						<IconChevronLeft />
						<span>Go back</span>
					</Box>
					{Object.entries(labelMap).map(([key, label]) => (
						<Tabs.Tab key={key} value={key}>
							{label}
						</Tabs.Tab>
					))}
				</div>
			</Tabs.List>

			<Tabs.Panel key={query} value={query}>
				<QueryTable query={query} data={data} />
			</Tabs.Panel>
		</Tabs>
	);
};
