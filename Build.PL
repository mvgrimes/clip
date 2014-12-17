use Module::Build;
my $build = Module::Build->new(
    module_name   => 'App::Clip',
    license       => 'perl',
    create_readme => 1,

    configure_requires => {
        'Module::Build' => '0.4004',
    },

    requires => {
        'perl' => '5.010',

        'App::Cmd'                  => '0.32',
        'AnyEvent'                  => 0,
        'AnyEvent::Filesys::Notify' => 0,
        'AnyEvent::Retry'           => '0.03',
        'Promises'                  => '0.04',
        'Data::Dump'                => 0,
        'IPC::Cmd'                  => 0,
        'Git::Wrapper'              => '0.030',
        'Sys::Hostname'             => 0,
        'Try::Tiny'                 => 0,

        'Moo'                                   => '1.001',
        'MooX::late'                            => '0.014',
        'Data::Perl'                            => '0.002007',
        'FindBin::libs'                         => '0',
        'Type::Tiny'                            => 0,
        'Types::Standard'                       => '0.001',
        'Log::Any'                              => '0.15',
        'Log::Any::Adapter'                     => '0.11',
        'Log::Any::Adapter::ScreenColoredLevel' => '0.04',
        'Config::General'                       => '2.52',
        'Growl::Tiny'                           => '0.0.4',
    },

    test_requires => {
        'Test::More'        => '0.87',
        'Test::Exception'   => '0.35',
        'Test::Class'       => '0.39',
        'Test::Differences' => '0.61',
        'Test::Most'        => '0.31',
    },

    meta_merge => {
        resources => {
            bugtracker => 'https://github.com/mvgrimes/perl-app-clip/issues',
            repository => 'https://github.com/mvgrimes/perl-app-clip'
        }
    },
);

$build->create_build_script;
