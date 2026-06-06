import { Modal, VisuallyHidden } from '@mantine/core';
import { DatePicker } from '@mantine/dates';
import { useDidUpdate, useMediaQuery } from '@mantine/hooks';
import { formatISO, sub } from 'date-fns';
import { useState } from 'react';
import { useSearchParams } from 'react-router';

import { Button, CloseButton } from '@/components/Button';

import classes from './DatePicker.module.css';

interface DatePickerRangeProps {
	open: boolean;
	setOpen: (open: boolean) => void;
}

type DateRangeValue = [string | null, string | null];

const today = () => formatISO(new Date(), { representation: 'date' });
const previousMonth = () =>
	formatISO(sub(new Date(), { months: 1 }), { representation: 'date' });

const DatePickerRange = ({ open, setOpen }: DatePickerRangeProps) => {
	const [searchParams, setSearchParams] = useSearchParams();

	const getDateRange = (): DateRangeValue => {
		const start = searchParams.get('start');
		const end = searchParams.get('end');
		if (searchParams.get('period') === 'custom' && start && end) {
			return [start, end];
		}
		return [null, null];
	};

	const [date, setDate] = useState<DateRangeValue>(getDateRange());

	useDidUpdate(() => {
		setDate(getDateRange());
	}, [searchParams]);

	const handleSubmit = () => {
		const [start, end] = date;
		if (start && end) {
			setSearchParams((params) => {
				params.set('period', 'custom');
				params.set('start', start);
				params.set('end', end);
				return params;
			});
		}
		setOpen(false);
	};

	const isMobile = useMediaQuery('(max-width: 48em)');

	return (
		<Modal
			opened={open}
			onClose={() => setOpen(false)}
			withCloseButton={false}
			centered
			size="auto"
			classNames={{
				overlay: classes.overlay,
				content: classes.content,
			}}
		>
			<div className={classes.header}>
				<h2 className={classes.title}>Select a date range</h2>
				<CloseButton label="Close date picker" onClick={() => setOpen(false)} />
			</div>
			<VisuallyHidden>
				Select the start and end date for the date range.
			</VisuallyHidden>
			<DatePicker
				type="range"
				className={classes.root}
				value={date}
				onChange={setDate}
				numberOfColumns={isMobile ? 1 : 2}
				defaultDate={previousMonth()}
				minDate="2024-01-01"
				maxDate={today()}
			/>
			<Button className={classes.apply} onClick={handleSubmit}>
				Apply
			</Button>
		</Modal>
	);
};

export { DatePickerRange };
