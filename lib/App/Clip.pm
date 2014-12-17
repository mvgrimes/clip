package App::Clip;

# ABSTRACT: clip

use 5.010;
use strict;
use warnings;

use App::Cmd::Setup -app;
use Log::Any::Adapter;

our $VERSION = '0.10';

Log::Any::Adapter->set( 'ScreenColoredLevel', min_level => 'debug' );

1;

__END__

=pod

=head1 NAME

DateTimeX::Immutable - An immutable subclass of DateTime

=head1 VERSION

version 0.33

=head1 SYNOPSIS

    $x

=head1 DESCRIPTION

Description

=head1 SEE ALSO

L<DateTime>

=head1 AUTHOR

Mark Grimes, E<lt>mgrimes@cpan.orgE<gt>

=head1 COPYRIGHT AND LICENSE

This software is copyright (c) 2014 by Mark Grimes, E<lt>mgrimes@cpan.orgE<gt>.

This is free software; you can redistribute it and/or modify it under
the same terms as the Perl 5 programming language system itself.

=cut
