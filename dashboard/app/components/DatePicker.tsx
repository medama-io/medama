import * as Dialog from '@radix-ui/react-dialog';
import { ChevronLeftIcon, ChevronRightIcon } from '@radix-ui/react-icons';
import * as VisuallyHidden from '@radix-ui/react-visually-hidden';
import { useSearchParams } from '@remix-run/react';
import { formatISO, parseISO, sub } from 'date-fns';
import { useState } from 'react';
import { type DateRange, DayPicker } from 'react-day-picker';

import { Button, CloseButton } from '@/components/Button';
import { useDidUpdate } from '@/hooks/use-did-update';
import { useMediaQuery } from '@/hooks/use-media-query';

import classes from './DatePicker.module.css';

interface DatePickerProps {
	open: boolean;
	setOpen: (open: boolean) => void;
}

const DatePickerRange = ({ open, setOpen }: DatePickerProps) => {
	const [searchParams, setSearchParams] = useSearchParams();

	const getDateRange = () => {
		const start = searchParams.get('start');
		const end = searchParams.get('end');
		if (searchParams.get('period') === 'custom' && start && end) {
			return {
				from: parseISO(start),
				to: parseISO(end),
			};
		}
		return undefined;
	};

	const [date, setDate] = useState<DateRange | undefined>(getDateRange());

	// If the search params change, we should verify if the date range has changed
	useDidUpdate(() => {
		setDate(getDateRange());
	}, [searchParams]);

	const handleSubmit = () => {
		if (date) {
			setSearchParams((params) => {
				params.set('period', 'custom');
				if (date.from)
					params.set('start', formatISO(date.from, { representation: 'date' }));
				if (date.to)
					params.set('end', formatISO(date.to, { representation: 'date' }));
				return params;
			});
		}
		setOpen(false);
	};

	const isMobile = useMediaQuery('(max-width: 48em)');

	return (
		<Dialog.Root open={open} onOpenChange={setOpen}>
			<Dialog.Portal>
				<Dialog.Overlay className={classes.overlay} />
				<Dialog.Content className={classes.content}>
					<div className={classes.header}>
						<Dialog.Title className={classes.title}>
							Select a date range
						</Dialog.Title>
						<Dialog.Close asChild>
							<CloseButton label="Close date picker" />
						</Dialog.Close>
					</div>
					<VisuallyHidden.Root asChild>
						<Dialog.Description>
							Select the start and end date for the date range.
						</Dialog.Description>
					</VisuallyHidden.Root>
					<DayPicker
						mode="range"
						classNames={classes}
						selected={date}
						onSelect={setDate}
						numberOfMonths={isMobile ? 1 : 2}
						components={{
							Chevron: (props) => {
								if (props.orientation === 'left') {
									return <ChevronLeftIcon {...props} />;
								}
								return <ChevronRightIcon {...props} />;
							},
						}}
						defaultMonth={sub(new Date(), { months: 1 })}
						startMonth={new Date(2024, 0)}
						endMonth={new Date()}
					/>
					<Button className={classes.apply} onClick={handleSubmit}>
						Apply
					</Button>
				</Dialog.Content>
			</Dialog.Portal>
		</Dialog.Root>
	);
};

export { classes as datePickerClasses, DatePickerRange };
